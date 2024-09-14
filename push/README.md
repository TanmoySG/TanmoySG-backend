# Updating Data for tanmoysg.com to wunderDb Collections using GitOps

This is the architectural design and thoughts behind this application.

## Preface

Personal website data is mostly unchanging and not dynamic. For this, data like social media accounts, profile details, etc need not be divided into domains and put into different databases or collections or tables.

We've migrated from the older data design where each domain had a collections, eg. different collections for social-media profiles, education, profile highlights, etc, to a newer, consolidated data design where data that wont change is clubbed into a single collection, i.e clubbed social-media profiles, education, profile highlights into a single schema and collection.

| Older Database Design              | New Database Design                |
| ---------------------------------- | ---------------------------------- |
| ![alt text](./old-data-design.png) | ![alt text](./new-data-design.png) |

| Collections (old)                   | Collections (new)                                                                                    |
|-------------------------------------|------------------------------------------------------------------------------------------------------|
| skills                              | [skills](../schema/databases/tsg-on-web_v0_beta_1/collections/skills/skills.schema.json)             |
| projects                            | [projects](../schema/databases/tsg-on-web_v0_beta_1/collections/projects/projects.schema.json)       |
| education, social, profileSpotlight | [profile](../schema/databases/tsg-on-web_v0_beta_1/collections/profile/profile.schema.json)          |
| messages, feedback                  | [messages](../schema/databases/tsg-on-web_v0_beta_1/collections/messages/messages.schema.json)       |
| experience                          | [experience](../schema/databases/tsg-on-web_v0_beta_1/collections/experience/experience.schema.json) |

Note how collections like `education`, `profileSpotlight`, `social` have been consolidated into a single schema `profile`, while `messages` and `feedback` have been merged into `messages`.

## Why Push Data using GitOps

The aforementioned data is often static, making them maintainable with JSON files. The data saved in these JSON files are not huge, so maintaining them as objects in a single file is also pretty easy.

That being said, for these data to be available via API calls, we need to have these pushed into [wunderDb](https://github.com/TanmoySG/wunderDB). With wdb maintaining and having data as json is very easy and effortless.

To make is easier, the idea is to remove the requirement of manually creating/ingesting new data or patching existing data when there is any change in the aforementioned JSON files. The only manual process should be to update the JSON files, rest of the things should be automated. Since we're making use of a git repository to maintain these files, it makes sense to use GitOps, i.e when file changes in the repository, GitHub should take care of the ingestion and updation automatically.

Read More about GitOps [here](https://about.gitlab.com/topics/gitops/).

## The GitOps Architecture

![alt text](data-push.drawio.png)

- When any record changes, a PR is raised with the change to the main branch. Each directory under `/data` corresponds to a collection.
- When the PR is merged to main, the [`data-sync`](../.github/workflows/data-sync.yaml) workflow runs. The workflow runs only if there is a change in the `/data` directory.
- The workflow runs an python script called [`push`](app.py).
- For each collection, i.e directory under `/data` , it gets the records already in the database.
- If the record already exists it patches the existing data, else create a new record.

Note: Since the records are reletively less, performing a patch operation even if nothing changed saves on field level value matching, which otherwise would be a compute intensitve process. This tradeoff is acceptable when there are less records.

## The Workflow

The [`data-sync.yaml`](../.github/workflows/data-sync.yaml) workflow configuration defines the steps to be run. This section talks about the steps involved.

```yaml
name: Sync Data

on:
  workflow_dispatch:
    inputs:
      confirm:
        type: boolean
        description: 'Confirm Manual Trigger for Sync Job'
        required: true
  push:
    branches: ['main']
    paths: ['data/*/records.json']
```

The first section of the YAML file contains the name of the workflow, in this case `Sync Data`. The second sections defines the triggers for running the workflow.

- The `push` trigger defines that the workflow should run when
  - There is a push in the `main` branch, defined in the `branches` field
  - And if there is a change in `/data/*/records.json` , i.e if there is a change in any of the records.json files under the directories within `data` directory, defined in the `paths` field.
  - The workflow won't run if BOTH the conditions are not met.
- In addition to the `push` trigger, sometimes the workflow needs to be run on-demand. The `workflow_dispatch` trigger is for running the workflow manually, i.e without any requirement of updating records.
  - It expects a confirmation input `confirm` of type boolean.
  - The workflow runs, when user wants to run is manually. It runs and updates/syncs the records in the `/data` directory.

```yaml
jobs: 
  sync: 
    runs-on: ubuntu-latest 
    steps: 
      - name: Checkout Repository
        uses: actions/checkout@v4
      - name: Set up Python 
        uses: actions/setup-python@v5
        with:
            python-version: '3.11'
```

The `jobs` section of the workflow file defines the jobs to run.

- The `runs-on` field under the job (`sync` in this case) defines the OS/Platform on which these steps should be run, called the `runner`.
- The first step `Checkout Repository` checksout or loads the repository onto the runner and uses the action[`actions/checkout@v4`](https://github.com/actions/checkout)
- The second step `Set up Python` uses the [`actions/setup-python@v5`](https://github.com/actions/setup-python) action to setup Python on the runner, the version of python to be installed can be specified in the `with.python-version` field.
- The actions - `actions/checkout@v4` and `actions/setup-python@v5` are available on GitHub marketplace.

```yaml
      - name: Build and Run 
        env: 
            BASE_URL: ${{ secrets.BASE_URL }}
            WDB_USERNAME: ${{ secrets.WDB_USERNAME }} 
            WDB_PASSWORD: ${{ secrets.WDB_PASSWORD }} 
        run: | 
            echo '### Sync Summary 📋' >> $GITHUB_STEP_SUMMARY
            trigger=$(echo ${{ github.event_name }})
            if [ $trigger == "workflow_dispatch" ]; then
              echo '💡 Trigger: `Manual`' >> $GITHUB_STEP_SUMMARY
            elif [ $trigger == "push" ]; then
              echo '💡 Trigger: `Record(s) Updated`' >> $GITHUB_STEP_SUMMARY
            fi
            cd ${GITHUB_WORKSPACE}/push
            pip install -r requirements.txt
            echo '```' >> $GITHUB_STEP_SUMMARY
            echo "🪵 Sync Run Logs" >> $GITHUB_STEP_SUMMARY
            echo >> $GITHUB_STEP_SUMMARY
            python3 app.py
            cat push.log 
            cat push.log >> $GITHUB_STEP_SUMMARY
            echo '```' >> $GITHUB_STEP_SUMMARY
            echo "✅ Sync Run Completed." >> $GITHUB_STEP_SUMMARY
```

This is the step where the actual application/python script runs. The step is called `Build and Run`, with a few environment variables - `BASE_URL` , `WDB_USERNAME` and `WDB_PASSWORD`, setup in the `env` section/field. The value of these are fetched from the ["Repository Secrets"](https://docs.github.com/en/actions/security-for-github-actions/security-guides/using-secrets-in-github-actions). The shell instructions to be run are setup/specified in the `run` section.

```shell
echo '### Sync Summary 📋' >> $GITHUB_STEP_SUMMARY
trigger=$(echo ${{ github.event_name }})
if [ $trigger == "workflow_dispatch" ]; then
    echo '💡 Trigger: `Manual`' >> $GITHUB_STEP_SUMMARY
elif [ $trigger == "push" ]; then
    echo '💡 Trigger: `Record(s) Updated`' >> $GITHUB_STEP_SUMMARY
fi
```

- These steps are pushing a markdown "summary" into the `STEP SUMMARY` section of the output. The `GITHUB_STEP_SUMMARY` contains the step output/summary and anything set/piped/appended to it is displayed in the Step Summary section in the GitHub UI.
- Different message is pushed to the step summary based on the trigger using the `$trigger == "workflow_dispatch"` checks.

```shell
cd ${GITHUB_WORKSPACE}/push
pip install -r requirements.txt
echo '```' >> $GITHUB_STEP_SUMMARY
echo "🪵 Sync Run Logs" >> $GITHUB_STEP_SUMMARY
echo >> $GITHUB_STEP_SUMMARY
python3 app.py
cat push.log
```

- In the runner, "we" move to the `/push` directory from the repsitory root directory, where the application is stored.
- The requirements/dependancies, stored in requirements.txt, are installed onto the runner using `pip`, followed by more step summary instructions.
- The application is run at the end using `python3 app.py`. The application pushes the logs into a log file on the runner - `push.log`.
- The logs are then displayed on the workflow run logs by `cat`-ing the log file.

```shell
cat push.log >> $GITHUB_STEP_SUMMARY
echo '```' >> $GITHUB_STEP_SUMMARY
echo "✅ Sync Run Completed." >> $GITHUB_STEP_SUMMARY
```

- The logs are also pushed into the step summary to be displayed in the GitHub Actions run Summary section.

![alt text](summary.png)

Read more about Step Summary [here](https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/workflow-commands-for-github-actions#adding-a-job-summary).

PS: I learnt about step summary while working on this workflow! 😄
