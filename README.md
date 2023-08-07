# auditr
Stores Microsoft Azure Activity Reports API data in a local database

### Requirements

# Azure

You will need to create an app registration in your azure portal, it will require the following permissions:
Microsoft Graph:
- AuditLog.Read.All
- Directory.Read.All
- User.Read

Once you've created the app registration with those permissions, you will need to record the tenant id, the 
client id and the client secret, and put these values into the `config/auditr.yaml` configuration file.

[This link has more information](https://learn.microsoft.com/en-us/graph/api/resources/azure-ad-auditlog-overview?view=graph-rest-1.0)

# Open AI / Chat GPT

You will need to create an API key for Chat GPT to enable db_chains natural language processing. These can be obtained at openai.com
and will require a credit card most likely. Once you've created an API key put the key into the `config/config.yaml`

### Initial setup:
This was developed on a Unbuntu 22.10, ymmv with other systems.

- Ensure you've done the steps above first.
 
- Clone this repository on a linux system:
`git clone git@github.com:t1nfoil/auditr.git`

- run:
  `./config.sh`, it will create the postgres database container + databases, and bring up the sqlgpt (inference) and auditr containers.

### Handy Commands
To get into the database, and view tables / data:
`psql -U postgres -h localhost`, default password is `password`

Connect to the database:
`\c audit;`

Drop the database:
`DROP DATABASE audit;`

View table's and schema:
`\d`

View docker containers:
`docker container ls`

Stop docker container:
`docker container stop <name>`

View docker logs:
`docker logs <container name>`

The sqlgpt gradio runs on `http://localhost:8080`, so that's where you want to point your browser to query with natural language, you can ask it:
show me the tables in the database, show me the columns in <tablename> etc.. 
