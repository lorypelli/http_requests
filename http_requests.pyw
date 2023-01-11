import PySimpleGUI
from requests import get, post, delete
from shutil import rmtree
from os import path
PySimpleGUI.theme("BlueMono")
def login():
    import_from_config = [
        [PySimpleGUI.Button("Import from config")]
    ]
    bot_token = [
        [PySimpleGUI.Text("Insert bot token"), PySimpleGUI.InputText(size=(100), key="tkn_textbox")],
        [PySimpleGUI.Button("Validate")]
    ]
    layout = [
        [import_from_config, bot_token]
    ]
    window = PySimpleGUI.Window("Login", layout, element_justification="c", icon="./app_icon.ico", font="Arial")
    while True:
        event, values = window.read()
        if (event == "Import from config"):
            try:
                import config
                if (path.exists("./__pycache__")):
                    rmtree("./__pycache__")
                response = get("https://discord.com/api/auth/login", headers = {
                    "authorization": "Bot " + config.bot_token
                })
                if (response.status_code == 200 and event == "Import from config"):
                    window["tkn_textbox"].Update(config.bot_token)
                elif (response.status_code != 200 and event == "Import from config"):
                    PySimpleGUI.popup("There was an error, try again!", no_titlebar=True)
            except:
                PySimpleGUI.popup("The file config.py doesn't exists or is not configured as it should be", no_titlebar=True)
        if (event == "Validate"):
            response = get("https://discord.com/api/auth/login", headers = {
                "authorization": "Bot " + values["tkn_textbox"]
            })
            if (response.status_code != 200 and event == "Validate"):
                PySimpleGUI.popup("There was an error, try again!", no_titlebar=True)
            elif (response.status_code == 200 and event == "Validate"):
                PySimpleGUI.popup("Validation Passed!", no_titlebar=True)
                login.tkn_value = values["tkn_textbox"]
                login.username = (get("https://discord.com/api/users/@me", headers = {
                    "authorization": "Bot " + values["tkn_textbox"]
                })).json()["username"]
                login.id = (get("https://discord.com/api/users/@me", headers = {
                    "authorization": "Bot " + values["tkn_textbox"]
                })).json()["id"]
                window.close()
                program()
        elif (event == PySimpleGUI.WIN_CLOSED):
            break
def program():
    logged_user = [
        [PySimpleGUI.Text("User ID " + login.id), PySimpleGUI.Push(), PySimpleGUI.Button("Logout"), PySimpleGUI.Push(), PySimpleGUI.Text("Logged in as " + login.username)]
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
        [PySimpleGUI.Text("Select action"), PySimpleGUI.Combo(["Write a message", "Delete a channel", "Delete a message"], size=(100), default_value="Write a message", readonly=True, key="selectbox"), PySimpleGUI.Button("Confirm")]
    ]
    message_id = [
        [PySimpleGUI.Text("Insert message id", key="msg_id_text", visible=False), PySimpleGUI.InputText(size=(100), key="msg_id_textbox", visible=False)]
    ]
    layout = [
        [logged_user, channel_id, list, message, message_id, send_delete_btn]
    ]
    window = PySimpleGUI.Window("http_requests", layout, element_justification="c", icon="./app_icon.ico", font="Arial")
    while True:
        event, values = window.read()
        if (event == "Validate"):
            response = get("https://discord.com/api/channels/" + values["chn_textbox"], headers = {
                "authorization": "Bot " + login.tkn_value
            })
            if (response.status_code != 200 and event == "Validate"):
                PySimpleGUI.popup("There was an error, try again!", no_titlebar=True)
            elif (response.status_code == 200 and event == "Validate"):
                PySimpleGUI.popup("Validation Passed!", no_titlebar=True)
        elif (event == "Logout"):
            response = PySimpleGUI.popup_ok_cancel("Are you sure you want to logout", no_titlebar=True)
            if (response == "OK"):
                window.close()
                login()
        elif (event == PySimpleGUI.WIN_CLOSED):
            break
        elif (values["selectbox"] == "Delete a channel"):
            window["msg_text"].Update(visible = False)
            window["msg_textbox"].Update(visible = False)
            window["msg_id_text"].Update(visible = False)
            window["msg_id_textbox"].Update(visible = False)
            window["Send_Delete"].Update("Delete") 
        elif (values["selectbox"] == "Delete a message"):
            window["msg_text"].Update(visible = False)
            window["msg_textbox"].Update(visible = False)
            window["msg_id_text"].Update(visible = True)
            window["msg_id_textbox"].Update(visible = True)
            window["Send_Delete"].Update("Delete") 
        else:
            window["msg_text"].Update(visible = True)
            window["msg_textbox"].Update(visible = True)
            window["msg_id_text"].Update(visible = False)
            window["msg_id_textbox"].Update(visible = False)
            window["Send_Delete"].Update("Send")
        if (values["selectbox"] != "Delete a channel" and values["selectbox"] != "Delete a message"):
            def post_msg(msg: str):
                response = post("https://discord.com/api/channels/" + values["chn_textbox"] + "/messages", headers = {
                    "authorization": "Bot " + login.tkn_value
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
                    "authorization": "Bot " + login.tkn_value
                })
                if (response.status_code != 200 and event == "Send_Delete"):
                    PySimpleGUI.popup("There was an error, try again!", no_titlebar=True)
                elif (response.status_code == 200 and event == "Send_Delete"):
                    PySimpleGUI.popup("The channel has been deleted successfully!", no_titlebar=True)
            if (event == "Send_Delete"):
                delete_chn(values["chn_textbox"])
        elif (values["selectbox"] == "Delete a message"):
            def delete_msg(msg_id: str):
                response = delete("https://discord.com/api/channels/" + values["chn_textbox"] + "/messages/" + msg_id, headers = {
                    "authorization": "Bot " + login.tkn_value
                })
                if (response.status_code != 204 and event == "Send_Delete"):
                    PySimpleGUI.popup("There was an error, try again!", no_titlebar=True)
                elif (response.status_code == 204 and event == "Send_Delete"):
                    PySimpleGUI.popup("The message has been deleted successfully!", no_titlebar=True)
            if (event == "Send_Delete"):
                delete_msg(values["msg_id_textbox"])
if (__name__ == "__main__"):
    login()