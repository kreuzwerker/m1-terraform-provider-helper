import requests 

RESULT_MARKDOWN = "docs/provider_information.md"
versions_available = 0

def getting_providers():
    print("starting provider search")
    provider_names = []
    page_count = 1

    while True:
        response = requests.get(f"https://registry.terraform.io/v2/providers?page[size]=100&page[number]={page_count}")

        yes = response.json()

        next_page = yes["meta"]["pagination"]["next-page"]
        total_pages = yes["meta"]["pagination"]["total-pages"]

        if page_count != total_pages:
            for i in range(len(yes["data"])):
                provider_names.append({"name": yes["data"][i]["attributes"]["full-name"], "github": yes["data"][i]["attributes"]["source"], 
                "terraform_link": f"https://registry.terraform.io/providers/{yes['data'][i]['attributes']['full-name']}"})
            page_count = next_page
        else:
            break
    
    print("done with the providers")
    return provider_names


def get_publishing_date_of_version(provider, version):
    res = requests.get(url=f"https://registry.terraform.io/v1/providers/{provider}/{version}")
    json = res.json()
    return json["published_at"].split("T")[0]

def get_darwin_arm64_information_of_provider(provider):
    global versions_available

    res = requests.get(url=f"https://registry.terraform.io/v1/providers/{provider}/versions")
    json = res.json()
    for version in json["versions"]:
        for platform in version["platforms"]:
            if platform["os"] == "darwin" and platform["arch"] == "arm64":
                published_at = get_publishing_date_of_version(provider, version['version'])
                versions_available += 1
                return  {"version": version["version"], "date": published_at}

    return {"version": "no version yet", "date": " "}

def get_versions(providers):
    print("starting version search")

    for index, provider in enumerate(providers):
        provider["version"] = get_darwin_arm64_information_of_provider(provider['name'])

        if index % 50 == 0:
            print(f"{index} versions found")

    print("done with the versions")
    return providers

def write_to_table(version_dict):

    print("beginning to write to table")

    global versions_available

    percentage = round(100/len(version_dict) * versions_available, 2)

    tablefile = open(RESULT_MARKDOWN, "w")
    tablefile.write("### Terraform providers supporting darwin arm64")

    tablefile = open(RESULT_MARKDOWN, "a")
    tablefile.write("\nThe following list gives an overview, which Terraform providers offer a version supporting darwin arm64 architecture.\nIf a version is available, the table provides the first published version supporting darwin arm64 and its publishing date.")
    tablefile.write(f"\n#### Current percentage of providers that offer a version supporting darwin arm64:")
    tablefile.write(f"\n## {percentage}%")
    tablefile.write("\nprovider | first version | publishing date | github repo | terraform registry")
    tablefile.write("\n --- | --- | --- | --- | ---")

    for i in version_dict:
        tablefile.write(f"\n{i['name']} | {i['version']['version']} | {i['version']['date']} | {i['github']} | {i['terraform_link']}")
    print("all done")
    tablefile.close()

providers = getting_providers()
versions = get_versions(providers)
write_to_table(versions)