import PySimpleGUI
from requests import get, post, delete, patch, put
from shutil import rmtree
from os import path
from importlib import reload
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
    window = PySimpleGUI.Window("Login", layout, element_justification="c", icon="./app_icon.ico", font="Arial", size=(850, 110))
    while True:
        event, values = window.read()
        if (event == "Import from config"):
            try:
                import config
                reload(config)
                configtkn = config.config()
                if (path.exists("./__pycache__")):
                    rmtree("./__pycache__")
                response = get("https://discord.com/api/auth/login", headers = {
                    "authorization": "Bot " + configtkn
                })
                if (response.status_code == 200 and event == "Import from config"):
                    window["tkn_textbox"].Update(configtkn)
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
    confirm_action_btn = [
        [PySimpleGUI.Button("Send", key="confirm_action_btn")]
    ]
    list = [
        [PySimpleGUI.Text("Select action"), PySimpleGUI.Combo(["Write a message", "Edit a message", "Pin a message", "Edit a channel", "Create a thread", "Delete a channel", "Delete a message", "Unpin a message"], size=(100), default_value="Write a message", readonly=True, key="selectbox"), PySimpleGUI.Button("Confirm")]
    ]
    message_id = [
        [PySimpleGUI.Text("Insert message id", key="msg_id_text", visible=False), PySimpleGUI.InputText(size=(100), key="msg_id_textbox", visible=False)],
        [PySimpleGUI.Button("Validate", visible=False, key="msg_id_btn")]
    ]
    channel_name = [
        [PySimpleGUI.Text("Insert channel name", key="chn_name_text", visible=False), PySimpleGUI.InputText(size=(100), key="chn_name_textbox", visible=False)]
    ]
    thread_name = [
        [PySimpleGUI.Text("Insert thread name", key="thread_name_text", visible=False), PySimpleGUI.InputText(size=(100), key="thread_name_textbox", visible=False)]
    ]
    layout = [
        [logged_user, channel_id, list, channel_name, message, message_id, thread_name, confirm_action_btn]
    ]
    window = PySimpleGUI.Window("http_requests", layout, element_justification="c", icon="./app_icon.ico", font="Arial", size=(1130, 420))
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
        elif (event == "msg_id_btn"):
            response = get("https://discord.com/api/channels/" + values["chn_textbox"] + "/messages/" + values["msg_id_textbox"], headers = {
                "authorization": "Bot " + login.tkn_value
            })
            if (response.status_code != 200 and event == "msg_id_btn"):
                PySimpleGUI.popup("There was an error, try again!", no_titlebar=True)
            elif (response.status_code == 200 and event == "msg_id_btn"):
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
            window["msg_id_btn"].Update(visible = False)
            window["chn_name_text"].Update(visible = False)
            window["chn_name_textbox"].Update(visible = False)
            window["thread_name_text"].Update(visible = False)
            window["thread_name_textbox"].Update(visible = False)
            window["confirm_action_btn"].Update("Delete")
        elif (values["selectbox"] == "Delete a message"):
            window["msg_id_text"].Update(visible = True)
            window["msg_id_textbox"].Update(visible = True)
            window["msg_id_btn"].Update(visible = True)
            window["msg_text"].Update(visible = False)
            window["msg_textbox"].Update(visible = False)
            window["chn_name_text"].Update(visible = False)
            window["chn_name_textbox"].Update(visible = False)
            window["thread_name_text"].Update(visible = False)
            window["thread_name_textbox"].Update(visible = False)
            window["confirm_action_btn"].Update("Delete")
        elif (values["selectbox"] == "Edit a message"):
            window["msg_id_text"].Update(visible = True)
            window["msg_id_textbox"].Update(visible = True)
            window["msg_id_btn"].Update(visible = True)
            window["msg_text"].Update(visible = True)
            window["msg_textbox"].Update(visible = True)
            window["chn_name_text"].Update(visible = False)
            window["chn_name_textbox"].Update(visible = False)
            window["thread_name_text"].Update(visible = False)
            window["thread_name_textbox"].Update(visible = False)
            window["confirm_action_btn"].Update("Edit")
        elif (values["selectbox"] == "Edit a channel"):
            window["chn_name_text"].Update(visible = True)
            window["chn_name_textbox"].Update(visible = True)
            window["msg_text"].Update(visible = False)
            window["msg_textbox"].Update(visible = False)
            window["msg_id_btn"].Update(visible = False)
            window["msg_id_text"].Update(visible = False)
            window["msg_id_textbox"].Update(visible = False)
            window["thread_name_text"].Update(visible = False)
            window["thread_name_textbox"].Update(visible = False)
            window["confirm_action_btn"].Update("Edit")
        elif (values["selectbox"] == "Create a thread"):
            window["thread_name_text"].Update(visible = True)
            window["thread_name_textbox"].Update(visible = True)
            window["msg_id_text"].Update(visible = True)
            window["msg_id_textbox"].Update(visible = True)
            window["msg_id_btn"].Update(visible = True)
            window["msg_text"].Update(visible = False)
            window["msg_textbox"].Update(visible = False)
            window["chn_name_text"].Update(visible = False)
            window["chn_name_textbox"].Update(visible = False)
            window["confirm_action_btn"].Update("Create")
        elif (values["selectbox"] == "Pin a message"):
            window["msg_id_text"].Update(visible = True)
            window["msg_id_textbox"].Update(visible = True)
            window["msg_id_btn"].Update(visible = True)
            window["msg_text"].Update(visible = False)
            window["msg_textbox"].Update(visible = False)
            window["chn_name_text"].Update(visible = False)
            window["chn_name_textbox"].Update(visible = False)
            window["thread_name_text"].Update(visible = False)
            window["thread_name_textbox"].Update(visible = False)
            window["confirm_action_btn"].Update("Pin")
        elif (values["selectbox"] == "Unpin a message"):
            window["msg_id_text"].Update(visible = True)
            window["msg_id_textbox"].Update(visible = True)
            window["msg_id_btn"].Update(visible = True)
            window["msg_text"].Update(visible = False)
            window["msg_textbox"].Update(visible = False)
            window["chn_name_text"].Update(visible = False)
            window["chn_name_textbox"].Update(visible = False)
            window["thread_name_text"].Update(visible = False)
            window["thread_name_textbox"].Update(visible = False)
            window["confirm_action_btn"].Update("Unpin")
        else:
            window["msg_text"].Update(visible = True)
            window["msg_textbox"].Update(visible = True)
            window["msg_id_text"].Update(visible = False)
            window["msg_id_textbox"].Update(visible = False)
            window["msg_id_btn"].Update(visible = False)
            window["chn_name_text"].Update(visible = False)
            window["chn_name_textbox"].Update(visible = False)
            window["thread_name_text"].Update(visible = False)
            window["thread_name_textbox"].Update(visible = False)
            window["confirm_action_btn"].Update("Send")
        if (values["selectbox"] != "Delete a channel" and values["selectbox"] != "Delete a message" and values["selectbox"] != "Edit a message" and values["selectbox"] != "Edit a channel" and values["selectbox"] != "Create a thread" and values["selectbox"] != "Pin a message" and values["selectbox"] != "Unpin a message"):
            def post_msg(msg: str):
                response = post("https://discord.com/api/channels/" + values["chn_textbox"] + "/messages", headers = {
                    "authorization": "Bot " + login.tkn_value
                }, json = {
                    "content": msg
                })
                if (response.status_code != 200 and event == "confirm_action_btn"):
                    PySimpleGUI.popup("There was an error, try again!", no_titlebar=True)
                elif (response.status_code == 200 and event == "confirm_action_btn"):
                    PySimpleGUI.popup("The message has been sent successfully!", no_titlebar=True)
            if (event == "confirm_action_btn"):
                post_msg(values["msg_textbox"])
        elif (values["selectbox"] == "Delete a channel"):
            def delete_chn(chn_id: str):
                response = delete("https://discord.com/api/channels/" + chn_id, headers = {
                    "authorization": "Bot " + login.tkn_value
                })
                if (response.status_code != 200 and event == "confirm_action_btn"):
                    PySimpleGUI.popup("There was an error, try again!", no_titlebar=True)
                elif (response.status_code == 200 and event == "confirm_action_btn"):
                    PySimpleGUI.popup("The channel has been deleted successfully!", no_titlebar=True)
            if (event == "confirm_action_btn"):
                delete_chn(values["chn_textbox"])
        elif (values["selectbox"] == "Delete a message"):
            def delete_msg(msg_id: str):
                response = delete("https://discord.com/api/channels/" + values["chn_textbox"] + "/messages/" + msg_id, headers = {
                    "authorization": "Bot " + login.tkn_value
                })
                if (response.status_code != 204 and event == "confirm_action_btn"):
                    PySimpleGUI.popup("There was an error, try again!", no_titlebar=True)
                elif (response.status_code == 204 and event == "confirm_action_btn"):
                    PySimpleGUI.popup("The message has been deleted successfully!", no_titlebar=True)
            if (event == "confirm_action_btn"):
                delete_msg(values["msg_id_textbox"])
        elif (values["selectbox"] == "Edit a message"):
            def edit_msg(msg_id: str, msg: str):
                response = patch("https://discord.com/api/channels/" + values["chn_textbox"] + "/messages/" + msg_id, headers = {
                    "authorization": "Bot " + login.tkn_value
                }, json = {
                    "content": msg
                })
                if (response.status_code != 200 and event == "confirm_action_btn"):
                    PySimpleGUI.popup("There was an error, try again!", no_titlebar=True)
                elif (response.status_code == 200 and event == "confirm_action_btn"):
                    PySimpleGUI.popup("The message has been edited successfully!", no_titlebar=True)
            if (event == "confirm_action_btn"):
                edit_msg(values["msg_id_textbox"], values["msg_textbox"])
        elif (values["selectbox"] == "Edit a channel"):
            def edit_chn(chn_id: str, chn_name: str):
                response = patch("https://discord.com/api/channels/" + chn_id, headers = {
                    "authorization": "Bot " + login.tkn_value
                }, json = {
                    "name": chn_name
                })
                if (response.status_code != 200 and event == "confirm_action_btn"):
                    PySimpleGUI.popup("There was an error, try again!", no_titlebar=True)
                elif (response.status_code == 200 and event == "confirm_action_btn"):
                    PySimpleGUI.popup("The channel has been edited successfully!", no_titlebar=True)
            if (event == "confirm_action_btn"):
                edit_chn(values["chn_textbox"], values["chn_name_textbox"])
        elif (values["selectbox"] == "Create a thread"):
            def create_thread(chn_id: str, msg_id: str, thread_name: str):
                response = post("https://discord.com/api/channels/" + chn_id + "/messages/" + msg_id + "/threads", headers = {
                    "authorization": "Bot " + login.tkn_value
                }, json = {
                    "name": thread_name
                })
                if (response.status_code != 201 and event == "confirm_action_btn"):
                    print(response.status_code)
                    PySimpleGUI.popup("There was an error, try again!", no_titlebar=True)
                elif (response.status_code == 201 and event == "confirm_action_btn"):
                    PySimpleGUI.popup("The thread has been created successfully!", no_titlebar=True)
            if (event == "confirm_action_btn"):
                create_thread(values["chn_textbox"], values["msg_id_textbox"], values["thread_name_textbox"])
        elif (values["selectbox"] == "Pin a message"):
            def pin_msg(msg_id: str):
                response = put("https://discord.com/api/channels/" + values["chn_textbox"] + "/pins/" + msg_id, headers = {
                    "authorization": "Bot " + login.tkn_value
                })
                if (response.status_code != 204 and event == "confirm_action_btn"):
                    PySimpleGUI.popup("There was an error, try again!", no_titlebar=True)
                elif (response.status_code == 204 and event == "confirm_action_btn"):
                    PySimpleGUI.popup("The message has been pinned successfully!", no_titlebar=True)
            if (event == "confirm_action_btn"):
                pin_msg(values["msg_id_textbox"])
        elif (values["selectbox"] == "Unpin a message"):
            def unpin_msg(msg_id: str):
                response = delete("https://discord.com/api/channels/" + values["chn_textbox"] + "/pins/" + msg_id, headers = {
                    "authorization": "Bot " + login.tkn_value
                })
                if (response.status_code != 204 and event == "confirm_action_btn"):
                    PySimpleGUI.popup("There was an error, try again!", no_titlebar=True)
                elif (response.status_code == 204 and event == "confirm_action_btn"):
                    PySimpleGUI.popup("The message has been unpinned successfully!", no_titlebar=True)
            if (event == "confirm_action_btn"):
                unpin_msg(values["msg_id_textbox"])
if (__name__ == "__main__"):
    login()