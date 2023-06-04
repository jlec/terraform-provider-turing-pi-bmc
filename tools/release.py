#!/usr/bin/env python
import hashlib
import os
import sys

import requests

try:
    from rich import print, print_json
    from rich.console import Console
    from rich.traceback import install

    install(show_locals=True)
    console = Console()
except ImportError:
    import json

    def print_json(data):
        print(json.dumps(data, indent=2))


GPG_FINGERPRINT = "FE007B07F3DCB6D4"


class TFERelease:
    organisation = "jlec"
    provider = "goer"
    namespace = "jlec"

    api_endpoint = f"https://app.terraform.io/api/v2/organizations/{organisation}"
    provider_endpoint = f"{api_endpoint}/registry-providers/private/{namespace}"

    provider_data = None
    version_data = None
    assets = None

    def __init__(self, version):
        self.version = version.strip()
        self.token = os.environ.get("TFE_TOKEN")
        self.header = {"Authorization": f"Bearer {self.token}", "Content-Type": "application/vnd.api+json"}

        self.create_provider()
        self.create_version()

        self.download_github_assets()
        self.upload_gpg()
        self.upload_shasums()
        self.upload_provider_platform()

    def get_provider(self) -> None:
        r = requests.get(f"{self.provider_endpoint}/{self.provider}", headers=self.header, timeout=10)
        if r.status_code == 200:
            print(f"Provider '{self.provider}' already created")
            self.provider_data = r.json()["data"]
            # print_json(data=self.provider_data)
            return

        print(f"Version '{self.version}' not released yet")
        # print_json(r.text)

    def create_provider(self) -> None:
        self.get_provider()
        if self.provider_data is not None:
            return
        payload = {
            "data": {
                "type": "registry-providers",
                "attributes": {"name": self.provider, "namespace": self.namespace, "registry-name": "private"},
            }
        }
        r = requests.post(f"{self.api_endpoint}/registry-providers", headers=self.header, json=payload, timeout=10)

        if r.status_code == 201:
            self.get_provider()
            # print_json(self.provider_data)
        else:
            print(f"Something went wrong when creating provider '{self.provider}'")
            print_json(data=r.json())
            sys.exit(1)

    def get_version(self):
        r = requests.get(f"{self.provider_endpoint}/{self.provider}/versions/{self.version}", headers=self.header, timeout=10)
        if r.status_code == 200:
            print(f"Version '{self.version}' already created")
            self.version_data = r.json()["data"]
            # print_json(data=self.version_data)
            return

        print(f"Version '{self.version}' not released yet")
        # print_json(data=j)

    def create_version(self):
        self.get_version()
        if self.version_data is not None:
            return
        payload = {
            "data": {
                "type": "registry-provider-versions",
                "attributes": {"version": self.version, "key-id": GPG_FINGERPRINT, "protocols": ["5.0"]},
            }
        }
        r = requests.post(f"{self.provider_endpoint}/{self.provider}/versions", headers=self.header, json=payload, timeout=10)

        if r.status_code == 201:
            self.version_data = r.json()["data"]
            # print_json(self.version_data)
        else:
            print(f"Something went wrong when creating version '{self.version}'")
            print_json(data=r.json())
            sys.exit(1)

    def download_github_assets(self):
        gh_api = f"https://api.github.com/repos/{self.organisation}/terraform-provider-{self.provider}"
        token = os.environ.get("GITHUB_TOKEN")
        header = {
            "Authorization": f"Bearer {token}",
            "Accept": "application/vnd.github+json",
            "X-GitHub-Api-Version": "2022-11-28",
        }
        # TODO: ERROR handling
        r = requests.get(f"{gh_api}/releases/latest", headers=header, timeout=10)
        self.assets = r.json()["assets"]
        # print_json(data=self.assets)

        header["Accept"] = "application/octet-stream"

        if not os.path.exists("releases"):
            os.mkdir("releases")

        for asset in self.assets:
            url = asset["url"]
            local_filename = f"releases/{asset['name']}"
            if os.path.exists(local_filename):
                continue

            print(f"Downloading file '{local_filename}'")
            with requests.get(url, headers=header, stream=True, timeout=10) as r:
                r.raise_for_status()
                with open(local_filename, "wb") as f:
                    for chunk in r.iter_content(chunk_size=8192):
                        # If you have chunk encoded response uncomment if
                        # and set chunk_size parameter to None.
                        # if chunk:
                        f.write(chunk)

    def upload_gpg(self):
        api_gpg = "https://app.terraform.io/api/registry/private/v2/gpg-keys"
        r = requests.get(f"{api_gpg}/{self.namespace}/{GPG_FINGERPRINT}", headers=self.header, timeout=10)
        if r.status_code == 404:
            with open(f"{GPG_FINGERPRINT}.asc", encoding="utf-8") as f:
                gpg_key = f.read()
            payload = {"data": {"type": "gpg-keys", "attributes": {"namespace": self.namespace, "ascii-armor": gpg_key}}}
            r = requests.post(api_gpg, headers=self.header, json=payload, timeout=10)
            if r.status_code != 200:
                print_json(data=r.json())
                sys.exit(1)

    def upload_shasums(self):
        shas = {
            "shasums-upload": f"terraform-provider-{self.provider}_{self.version}_SHA256SUMS",
            "shasums-sig-upload": f"terraform-provider-{self.provider}_{self.version}_SHA256SUMS.sig",
        }
        for k, v in shas.items():
            if k in self.version_data["links"]:
                with open(f"releases/{v}", "rb") as f:
                    r = requests.put(self.version_data["links"][k], files={v: f}, timeout=10)
                    r.raise_for_status()
                    self.get_version()
        # print_json(data=self.version_data)

    def get_provider_platform(self, os, arch) -> dict | None:
        r = requests.get(
            f"{self.provider_endpoint}/{self.provider}/versions/{self.version}/platforms/{os}/{arch}",
            headers=self.header,
            timeout=10,
        )
        if r.status_code == 200:
            return r.json()
        else:
            return None

    def create_provider_platform(self, os, arch, filename) -> dict:
        provider_platform = self.get_provider_platform(os, arch)
        if provider_platform is None:
            print(f"Creating platform for '{os}/{arch}'")
            with open(f"releases/{filename}", "rb") as f:
                data = f.read()
                sha256hash = hashlib.sha256(data).hexdigest()

            payload = {
                "data": {
                    "type": "registry-provider-version-platforms",
                    "attributes": {"os": os, "arch": arch, "shasum": sha256hash, "filename": filename},
                }
            }
            r = requests.post(
                f"{self.provider_endpoint}/{self.provider}/versions/{self.version}/platforms",
                headers=self.header,
                json=payload,
                timeout=10,
            )
            r.raise_for_status()

        provider_platform = self.get_provider_platform(os, arch)

        return provider_platform["data"]

    def upload_provider_platform(self):
        for asset in self.assets:
            filename = asset["name"]
            split_filename = filename.split("_")
            if len(split_filename) < 4:
                continue
            os = split_filename[2]
            arch = split_filename[3].rstrip(".zip")

            pp = self.create_provider_platform(os, arch, filename)
            # print_json(data=pp)

            if "provider-binary-upload" in pp["links"]:
                print(f"Uploading '{filename}' to TFC")
                with open(f"releases/{filename}", "rb") as f:
                    r = requests.put(pp["links"]["provider-binary-upload"], files={filename: f}, timeout=10)
                    r.raise_for_status()


if __name__ == "__main__":
    cwd = os.path.abspath(os.path.dirname(__file__))
    with open(f"{cwd}/../VERSION", encoding="utf-8") as f:
        provider_version = f.read()

    tfe = TFERelease(provider_version)
