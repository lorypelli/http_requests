import PySimpleGUI
from requests import get, post, delete
from shutil import rmtree
from os import path
import config
if (path.exists("./__pycache__")):
    rmtree("./__pycache__")
PySimpleGUI.theme("BlueMono")
logged_in_as = [
    [PySimpleGUI.Push(), PySimpleGUI.Text("To login, provide a valid bot token", key="logged_in_as")]
]
import_from_config = [
    [PySimpleGUI.Button("Import from config")]
]
bot_token = [
    [PySimpleGUI.Text("Insert bot token"), PySimpleGUI.InputText(size=(100), key="tkn_textbox")],
    [PySimpleGUI.Button("Validate")]
]
channel_id = [
    [PySimpleGUI.Text("Insert channel id"), PySimpleGUI.InputText(size=(100), key="chn_textbox")],
    [PySimpleGUI.Button("Validate")]
]
message = [
    [PySimpleGUI.Text("Insert message", key="msg_text"), PySimpleGUI.Multiline(size=(100, 5), key="msg_textbox")]
]
send_delete_btn = [
    [PySimpleGUI.Button("Send", key="Send_Delete")]
]
list = [
    [PySimpleGUI.Text("Select action"), PySimpleGUI.Combo(["Write a message", "Delete a channel"], size=(100), default_value="Write a message", readonly=True, key="selectbox"), PySimpleGUI.Button("Confirm")]
]
layout = [
    [logged_in_as, import_from_config, bot_token, channel_id, list, message, send_delete_btn]
]
window = PySimpleGUI.Window("http_requests", layout, element_justification="c", icon="./app_icon.ico", font="Arial")
while True:
    event, values = window.read()
    if (event == "Import from config"):
        response = get("https://discord.com/api/auth/login", headers = {
            "authorization": "Bot " + config.bot_token
        })
        if (response.status_code == 200):
            window["tkn_textbox"].Update(config.bot_token)
        response = get("https://discord.com/api/channels/" + config.channel_id, headers = {
            "authorization": "Bot " + config.bot_token
        })
        if (response.status_code == 200):
            window["chn_textbox"].Update(config.channel_id)
    if (event == "Validate"):
        response = get("https://discord.com/api/auth/login", headers = {
            "authorization": "Bot " + values["tkn_textbox"]
        })
        if (response.status_code != 200 and event == "Validate"):
            PySimpleGUI.popup("There was an error, try again!", no_titlebar=True)
            window["logged_in_as"].Update("To login, provide a valid bot token")
        elif (response.status_code == 200 and event == "Validate"):
            PySimpleGUI.popup("Validation Passed!", no_titlebar=True)
            response = get("https://discord.com/api/users/@me", headers = {
                "authorization": "Bot " + values["tkn_textbox"]
            })
            window["logged_in_as"].Update("Logged in as " + response.json()["username"])
    elif (event == "Validate0"):
        response = get("https://discord.com/api/channels/" + values["chn_textbox"], headers = {
            "authorization": "Bot " + values["tkn_textbox"]
        })
        if (response.status_code != 200 and event == "Validate0"):
            PySimpleGUI.popup("There was an error, try again!", no_titlebar=True)
        elif (response.status_code == 200 and event == "Validate0"):
            PySimpleGUI.popup("Validation Passed!", no_titlebar=True)
    elif (event == PySimpleGUI.WIN_CLOSED):
        break
    elif (values["selectbox"] == "Delete a channel"):
        window["msg_text"].Update(visible = False)
        window["msg_textbox"].Update(visible = False)
        window["Send_Delete"].Update("Delete")
    else:
        window["msg_text"].Update(visible = True)
        window["msg_textbox"].Update(visible = True)
        window["Send_Delete"].Update("Send")
    if (values["selectbox"] != "Delete a channel"):
        def post_msg(msg: str):
            response = post("https://discord.com/api/channels/" + values["chn_textbox"] + "/messages", headers = {
                "authorization": "Bot " + values["tkn_textbox"]
            }, data = {
                "content": msg
            })
            if (response.status_code != 200 and event == "Send_Delete"):
                PySimpleGUI.popup("There was an error, try again!", no_titlebar=True)
            elif (response.status_code == 200 and event == "Send_Delete"):
                PySimpleGUI.popup("The message has been sent successfully!", no_titlebar=True)
        if (event == "Send_Delete"):
            post_msg(values["msg_textbox"])
    elif (values["selectbox"] == "Delete a channel"):
        def delete_chn(chn_id: str):
            response = delete("https://discord.com/api/channels/" + chn_id, headers = {
                "authorization": "Bot " + values["tkn_textbox"]
            })
            if (response.status_code != 200 and event == "Send_Delete"):
                PySimpleGUI.popup("There was an error, try again!", no_titlebar=True)
            elif (response.status_code == 200 and event == "Send_Delete"):
                PySimpleGUI.popup("The channel has been deleted successfully!", no_titlebar=True)
        if (event == "Send_Delete"):
            delete_chn(values["chn_textbox"])