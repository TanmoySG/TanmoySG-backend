# Updating Data for tanmoysg.com to wunderDb Collections using GitOps

This is the architectural and thoughts behind this.

## Preface

Personal website data is mostly unchanging and not dynamic. For this, data like social media accounts, profile details, etc need not be divided into domains and put into different databases or collections or tables.

We've migrated from the older data design where each domain had a collections, eg. different collections for social-media profiles, education, profile highlights, etc, to a newer, consolidated data design where data that wont change is clubbed into a single collection, i.e clubbed social-media profiles, education, profile highlights into a single schema and collection.

Older Database Design
![alt text](./old-data-design.png)

New Database Design
![alt text](./new-data-design.png)

Note how collections like `education`, `profileSpotlight`, `social` have been moved to a single collection [`profile`](../schema/databases/tsg-on-web_v0_beta_1/collections/profile/profile.sample.json).

## Why Push Data using GitOps

The aforementioned data is often static, making them maintainable with JSON files. The data saved in these JSON files are not huge, so maintaining them as objects in a single file is also pretty easy.

That being said, for these data to be available via API calls, we need to have these pushed into [wunderDb](https://github.com/TanmoySG/wunderDB). With wdb maintaining and having data as json is very easy and effortless.

To make is easier, the idea is to remove the requirement of manually creating/ingesting new data or patching existing data when there is any change in the aforementioned JSON files. The only manual process should be to update the JSON files, rest of the things should be automated. Since we're making use of a git repository to maintain these files, it makes sense to use GitOps, i.e when file changes in the repository, GitHub should take care of the ingestion and updation automatically. 

Read More about GitOps [here](https://about.gitlab.com/topics/gitops/).

## The GitOps Architecture

![alt text](data-push.drawio.png)