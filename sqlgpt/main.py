import os
import yaml
import openai
import gradio as gr
from sqlalchemy import create_engine, MetaData,  select, Table, Column, String
from llama_index import LLMPredictor, ServiceContext, SQLDatabase, VectorStoreIndex, GPTSQLStructStoreIndex
from llama_index.indices.struct_store import SQLTableRetrieverQueryEngine
from llama_index.objects import SQLTableNodeMapping, ObjectIndex, SQLTableSchema
from langchain import OpenAI, SQLDatabaseChain
from langchain.chat_models import ChatOpenAI

# Load the configuration file
with open("/config/config.yaml", "r") as yamlfile:
    cfg = yaml.safe_load(yamlfile)

# If cfg["apikey"] is not set, exit with error apikey not set
if cfg["apikey"] == "":
    print("Error: apikey not set in config.yaml")
    exit()

os.environ["OPENAI_API_KEY"] = cfg["apikey"]
openai.api_key = cfg["apikey"]

print(cfg["apikey"])
#connect to database
pg_uri = f"postgresql+psycopg2://postgres:password@postgres-dev:5432/audit"
db = SQLDatabase.from_uri(pg_uri)

db_chain = SQLDatabaseChain(llm=ChatOpenAI(temperature=0, model_name='gpt-3.5-turbo-16k'), database=db, verbose=True)


def chatbot(input_text):
    output = db_chain.run(input_text)
    return output

iface = gr.Interface(fn=chatbot,
    inputs=gr.inputs.Textbox(lines=7, label="Enter your text"),
    outputs=[
        gr.outputs.Textbox(label="Result"),
    ],
    title="SecurityGPT")

iface.launch(server_name="0.0.0.0", server_port=8080)
