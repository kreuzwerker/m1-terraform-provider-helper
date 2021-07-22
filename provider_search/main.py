import requests 

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


def get_versions(clean_provider_names):

    print("starting version search")

    c = 0

    for i in clean_provider_names:

        response = requests.get(url=f"https://registry.terraform.io/v1/providers/{i['name']}/versions")

        yes2 = response.json()

            
        def getting_versions():

            global versions_available

            try:
                for m in yes2["versions"]:
                    for h in m["platforms"]:
                        if h["os"] == "darwin" and h["arch"] == "arm64":
                            next_response = requests.get(url=f"https://registry.terraform.io/v1/providers/{i['name']}/{m['version']}")
                            next_yes = next_response.json()
                            i["version"] = {"version": m["version"], "date": next_yes["published_at"].split("T")[0]}
                            versions_available += 1
                            return 
                i["version"] = {"version": "no version yet", "date": " "}
                return 
            except:
                print(f"{i} hello")
                return  
        
        c += 1

        if c % 50 == 0:
            print(f"{c} versions found")
        
        getting_versions()
    #print(clean_provider_names)
    print("done with the versions")
    return clean_provider_names



def write_to_table(version_dict):

    print("beginning to write to table")

    global versions_available

    percentage = round(100/len(version_dict) * versions_available, 2)

    tablefile = open("real_table.md", "w")
    tablefile.write("### Terraform providers supporting darwin arm64")

    tablefile = open("real_table.md", "a")
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





