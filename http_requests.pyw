from customtkinter import set_appearance_mode, set_default_color_theme, CTkLabel, CTkEntry, CTkButton, CTk, CTkComboBox, StringVar, CTkTextbox
from tkinter import messagebox
from requests import get, post, delete, patch, put
from webbrowser import open
from urllib import request
from os import environ, getenv, path, makedirs
from asyncio import run
from websockets.sync.client import connect
from json import dumps
user_version = "EighteenthRelease"
is_alpha = False
icon_url = "https://raw.githubusercontent.com/LoryPelli/http_requests/main/app_icon.ico"
icon_directory = f"{environ.get('SystemDrive')}/Users/{getenv('Username')}/http_requests"
icon_file = f"{icon_directory}/app_icon.ico"
set_appearance_mode("System")
set_default_color_theme("dark-blue")
if (not(path.exists(icon_directory))):
    makedirs(icon_directory)
try:
    request.urlretrieve(icon_url, icon_file)
except:
    pass
async def connect_to_gateway(token: str):
    ws = connect("wss://gateway.discord.gg/?v=10&encoding=json")
    payload = {
        "op": 2,
        "d": {
            "token": token,
            "intents": 0,
            "properties": {
                "os": "linux",
                "browser": "http_requests",
                "device": "discord"
            },
            "presence": {
                "activities": [{
                    "name": "http_requests",
                    "type": 3
                }],
                "status": "dnd"
            }
        }
    }
    ws.send(dumps(payload))
def login():
    app = CTk()
    app.title("Login")
    app.resizable(False, False)
    app.geometry("840x70+500+500")
    try:
        app.wm_iconbitmap(icon_file)
    except:
        pass
    def loginbtn():
        async def connection():
            await connect_to_gateway(tkn_textbox.get())
        try:
            response = post("https://discord.com/api/v10/auth/login", headers = {
                "authorization": f"Bot {tkn_textbox.get()}"
            })
        except:
            messagebox.showerror("Error", "Check your internet connection")
        if (response.status_code != 200):
            messagebox.showerror("Error", response.json()["message"])
        elif (response.status_code == 200):
            run(connection())
            messagebox.showinfo("Success", "Validation Passed!")
            login.tkn_value = tkn_textbox.get()
            login.username = (get("https://discord.com/api/v10/users/@me", headers = {
                "authorization": f"Bot {tkn_textbox.get()}"
            })).json()["username"]
            login.id = (get("https://discord.com/api/v10/users/@me", headers = {
                "authorization": f"Bot {tkn_textbox.get()}"
            })).json()["id"]
            app.destroy()
            program().mainloop()
    CTkLabel(app, text="Insert bot token", font=("Arial", 16)).place(relx=0.01, rely=0.1)
    tkn_textbox = CTkEntry(app, width=700, height=25, font=("Arial", 16))
    tkn_textbox.place(relx=0.15, rely=0.1)
    CTkButton(app, text="Validate", command=loginbtn, font=("Arial", 16)).place(relx=0.4, rely=0.5)
    return app
def program():
    app = CTk()
    app.title("http_requests")
    app.resizable(False, False)
    app.geometry("500x250+450+450")
    try:
        app.wm_iconbitmap(icon_file)
    except:
        pass
    def logout():
        response = messagebox.askyesno("Logout", "Are you sure you want to logout")
        if (response == True):
            async def disconnect(token: str):
                ws = connect("wss://gateway.discord.gg/?v=10&encoding=json")
                payload = {
                    "op": 2,
                    "d": {
                        "token": token,
                        "intents": 0,
                        "properties": {
                            "os": "linux",
                            "browser": "http_requests",
                            "device": "discord"
                        },
                        "presence": {
                            "activities": [{
                                "name": "disconnecting...",
                                "type": 3
                            }],
                            "status": "idle"
                        }
                    }
                }
                ws.send(dumps(payload))
            async def disconnection():
                await disconnect(login.tkn_value)
            run(disconnection())
            app.destroy()
            login().mainloop()
    app.protocol("WM_DELETE_WINDOW", logout)
    def combochoice(choice: str):
        if (choice == "Write a message"):
            msg_label.place(relx=0.01, rely=0.35)
            msg_textbox.place(relx=0.3, rely=0.35)
            msg_textbox.configure(height=100)
            chn_id_label.place(relx=0.01, rely=0.1)
            chn_id_textbox.place(relx=0.3, rely=0.1)
            chn_name_label.place_forget()
            chn_name_textbox.place_forget()
            msg_id_label.place_forget()
            msg_id_textbox.place_forget()
            thread_name_label.place_forget()
            thread_name_textbox.place_forget()
            guild_id_label.place_forget()
            guild_id_textbox.place_forget()
            chn_type_label.place_forget()
            chn_type.place_forget()
            usr_id_label.place_forget()
            usr_id_textbox.place_forget()
            confirm_action.configure(text="Send")
        elif (choice == "Edit a channel"):
            chn_name_label.place(relx=0.01, rely=0.35)
            chn_name_textbox.place(relx=0.3, rely=0.35)
            chn_id_label.place(relx=0.01, rely=0.1)
            chn_id_textbox.place(relx=0.3, rely=0.1)
            msg_label.place_forget()
            msg_textbox.place_forget()
            msg_id_label.place_forget()
            msg_id_textbox.place_forget()
            thread_name_label.place_forget()
            thread_name_textbox.place_forget()
            guild_id_label.place_forget()
            guild_id_textbox.place_forget()
            chn_type_label.place_forget()
            chn_type.place_forget()
            usr_id_label.place_forget()
            usr_id_textbox.place_forget()
            confirm_action.configure(text="Edit")
        elif (choice == "Edit a message"):
            msg_id_label.place(relx=0.01, rely=0.35)
            msg_id_textbox.place(relx=0.3, rely=0.35)
            msg_label.place(relx=0.01, rely=0.47)
            msg_textbox.place(relx=0.3, rely=0.47)
            msg_textbox.configure(height=75)
            chn_id_label.place(relx=0.01, rely=0.1)
            chn_id_textbox.place(relx=0.3, rely=0.1)
            chn_name_label.place_forget()
            chn_name_textbox.place_forget()
            thread_name_label.place_forget()
            thread_name_textbox.place_forget()
            guild_id_label.place_forget()
            guild_id_textbox.place_forget()
            chn_type_label.place_forget()
            chn_type.place_forget()
            usr_id_label.place_forget()
            usr_id_textbox.place_forget()
            confirm_action.configure(text="Edit")
        elif (choice == "Pin a message"):
            msg_id_label.place(relx=0.01, rely=0.35)
            msg_id_textbox.place(relx=0.3, rely=0.35)
            chn_id_label.place(relx=0.01, rely=0.1)
            chn_id_textbox.place(relx=0.3, rely=0.1)
            msg_label.place_forget()
            msg_textbox.place_forget()
            chn_name_label.place_forget()
            chn_name_textbox.place_forget()
            thread_name_label.place_forget()
            thread_name_textbox.place_forget()
            guild_id_label.place_forget()
            guild_id_textbox.place_forget()
            chn_type_label.place_forget()
            chn_type.place_forget()
            usr_id_label.place_forget()
            usr_id_textbox.place_forget()
            confirm_action.configure(text="Pin")
        elif (choice in ["Delete a message", "Unpin a message"]):
            msg_id_label.place(relx=0.01, rely=0.35)
            msg_id_textbox.place(relx=0.3, rely=0.35)
            chn_id_label.place(relx=0.01, rely=0.1)
            chn_id_textbox.place(relx=0.3, rely=0.1)
            msg_label.place_forget()
            msg_textbox.place_forget()
            chn_name_label.place_forget()
            chn_name_textbox.place_forget()
            thread_name_label.place_forget()
            thread_name_textbox.place_forget()
            guild_id_label.place_forget()
            guild_id_textbox.place_forget()
            chn_type_label.place_forget()
            chn_type.place_forget()
            usr_id_label.place_forget()
            usr_id_textbox.place_forget()
            if (choice == "Delete a message"):
                confirm_action.configure(text="Delete")
            elif (choice == "Unpin a message"):
                confirm_action.configure(text="Unpin")
        elif (choice == "Delete a channel"):
            chn_id_label.place(relx=0.01, rely=0.1)
            chn_id_textbox.place(relx=0.3, rely=0.1)
            msg_id_label.place_forget()
            msg_id_textbox.place_forget()
            msg_label.place_forget()
            msg_textbox.place_forget()
            chn_name_label.place_forget()
            chn_name_textbox.place_forget()
            thread_name_label.place_forget()
            thread_name_textbox.place_forget()
            guild_id_label.place_forget()
            guild_id_textbox.place_forget()
            chn_type_label.place_forget()
            chn_type.place_forget()
            usr_id_label.place_forget()
            usr_id_textbox.place_forget()
            confirm_action.configure(text="Delete")
        elif (choice == "Create a thread"):
            msg_id_label.place(relx=0.01, rely=0.35)
            msg_id_textbox.place(relx=0.3, rely=0.35)
            thread_name_label.place(relx=0.01, rely=0.465)
            thread_name_textbox.place(relx=0.3, rely=0.465)
            chn_id_label.place(relx=0.01, rely=0.1)
            chn_id_textbox.place(relx=0.3, rely=0.1)
            msg_label.place_forget()
            msg_textbox.place_forget()
            chn_name_label.place_forget()
            chn_name_textbox.place_forget()
            guild_id_label.place_forget()
            guild_id_textbox.place_forget()
            chn_type_label.place_forget()
            chn_type.place_forget()
            usr_id_label.place_forget()
            usr_id_textbox.place_forget()
            confirm_action.configure(text="Create")
        elif (choice == "Create a channel"):
            guild_id_label.place(relx=0.01, rely=0.1)
            guild_id_textbox.place(relx=0.3, rely=0.1)
            chn_type_label.place(relx=0.01, rely=0.35)
            chn_type.place(relx=0.3, rely=0.35)
            chn_name_label.place(relx=0.01, rely=0.465)
            chn_name_textbox.place(relx=0.3, rely=0.480)
            msg_id_label.place_forget()
            msg_id_textbox.place_forget()
            thread_name_label.place_forget()
            thread_name_textbox.place_forget()
            msg_label.place_forget()
            msg_textbox.place_forget()
            chn_id_label.place_forget()
            chn_id_textbox.place_forget()
            usr_id_label.place_forget()
            usr_id_textbox.place_forget()
            confirm_action.configure(text="Create")
        elif (choice in ["Kick a user", "Ban a user", "Unban a user"]):
            guild_id_label.place(relx=0.01, rely=0.1)
            guild_id_textbox.place(relx=0.3, rely=0.1)
            usr_id_label.place(relx=0.01, rely=0.35)
            usr_id_textbox.place(relx=0.3, rely=0.35)
            msg_id_label.place_forget()
            msg_id_textbox.place_forget()
            msg_label.place_forget()
            msg_textbox.place_forget()
            chn_name_label.place_forget()
            chn_name_textbox.place_forget()
            thread_name_label.place_forget()
            thread_name_textbox.place_forget()
            chn_type_label.place_forget()
            chn_type.place_forget()
            chn_id_label.place_forget()
            chn_id_textbox.place_forget()
            if (choice == "Kick a user"):
                confirm_action.configure(text="Kick")
            elif (choice == "Ban a user"):
                confirm_action.configure(text="Ban")
            elif (choice == "Unban a user"):
                confirm_action.configure(text="Unban")
    def confirm():
        async def connection():
            await connect_to_gateway(login.tkn_value)
        choice = combobox.get()
        def post_msg(msg: str):
            try:
                response = post(f"https://discord.com/api/v10/channels/{chn_id_textbox.get()}/messages", headers = {
                    "authorization": f"Bot {login.tkn_value}"
                }, json = {
                    "content": msg
                })
            except:
                messagebox.showerror("Error", "Check your internet connection")
            if (response.status_code != 200):
                messagebox.showerror("Error", response.json()["message"])
            elif (response.status_code == 200):
                run(connection())
                messagebox.showinfo("Success", "The message has been successfully sent!")
        def edit_msg(msg_id: str, msg: str):
            response = patch(f"https://discord.com/api/v10/channels/{chn_id_textbox.get()}/messages/{msg_id}", headers = {
                "authorization": f"Bot {login.tkn_value}"
            }, json = {
                "content": msg
            })
            if (response.status_code != 200):
                messagebox.showerror("Error", response.json()["message"])
            elif (response.status_code == 200):
                run(connection())
                messagebox.showinfo("Success", "The message has been successfully edited!")
        def delete_msg(msg_id: str):
            response = delete(f"https://discord.com/api/v10/channels/{chn_id_textbox.get()}/messages/{msg_id}", headers = {
                "authorization": f"Bot {login.tkn_value}"
            })
            if (response.status_code != 204):
                messagebox.showerror("Error", response.json()["message"])
            elif (response.status_code == 204):
                run(connection())
                messagebox.showinfo("Success", "The message has been successfully deleted!")
        def pin_msg(msg_id: str):
            try:
                response = put(f"https://discord.com/api/v10/channels/{chn_id_textbox.get()}/pins/{msg_id}", headers = {
                    "authorization": f"Bot {login.tkn_value}"
                })
            except:
                messagebox.showerror("Error", "Check your internet connection")
            if (response.status_code != 204):
                messagebox.showerror("Error", response.json()["message"])
            elif (response.status_code == 204):
                run(connection())
                messagebox.showinfo("Success", "The message has been successfully pinned!")
        def unpin_msg(msg_id: str):
            try:
                response = delete(f"https://discord.com/api/v10/channels/{chn_id_textbox.get()}/pins/{msg_id}", headers = {
                    "authorization": f"Bot {login.tkn_value}"
                })
            except:
                messagebox.showerror("Error", "Check your internet connection")
            if (response.status_code != 204):
                messagebox.showerror("Error", response.json()["message"])
            elif (response.status_code == 204):
                run(connection())
                messagebox.showinfo("Success", "The message has been successfully unpinned!")
        def create_chn(guild_id: str, chn_name: str, choice: str):
            if (choice == "Text"):
                choice = 0
            elif (choice == "Voice"):
                choice = 2
            elif (choice == "Announcement"):
                choice = 5
            elif (choice == "Stage"):
                choice = 13
            elif (choice == "Forum"):
                choice = 15
            try:
                response = post(f"https://discord.com/api/v10/guilds/{guild_id}/channels", headers = {
                    "authorization": f"Bot {login.tkn_value}"
                }, json = {
                    "name": chn_name,
                    "type": choice
                })
            except:
                messagebox.showerror("Error", "Check your internet connection")
            if (response.status_code != 201):
                messagebox.showerror("Error", response.json()["message"])
            elif (response.status_code == 201):
                run(connection())
                messagebox.showinfo("Success", "The channel has been successfully created!")
        def edit_chn(chn_id: str, chn_name: str):
            try:
                response = patch(f"https://discord.com/api/v10/channels/{chn_id}", headers = {
                    "authorization": f"Bot {login.tkn_value}"
                }, json = {
                    "name": chn_name
                })
            except:
                messagebox.showerror("Error", "Check your internet connection")
            if (response.status_code != 200):
                messagebox.showerror("Error", response.json()["message"])
            elif (response.status_code == 200):
                run(connection())
                messagebox.showinfo("Success", "The channel has been successfully edited!")
        def delete_chn(chn_id: str):
            try:
                response = delete(f"https://discord.com/api/v10/channels/{chn_id}", headers = {
                    "authorization": f"Bot {login.tkn_value}"
                })
            except:
                messagebox.showerror("Error", "Check your internet connection")
            if (response.status_code != 200):
                messagebox.showerror("Error", response.json()["message"])
            elif (response.status_code == 200):
                run(connection())
                messagebox.showinfo("Success", "The channel has been successfully deleted!")
        def create_thread(msg_id: str, thread_name: str):
            try:
                response = post(f"https://discord.com/api/v10/channels/{chn_id_textbox.get()}/messages/{msg_id}/threads", headers = {
                    "authorization": f"Bot {login.tkn_value}"
                }, json = {
                    "name": thread_name
                })
            except:
                messagebox.showerror("Error", "Check your internet connection")
            if (response.status_code != 201):
                messagebox.showerror("Error", response.json()["message"])
            elif (response.status_code == 201):
                run(connection())
                messagebox.showinfo("Success", "The thread has been successfully created!")
        def kick_usr(guild_id: str, usr_id: str):
            try:
                response = delete(f"https://discord.com/api/v10/guilds/{guild_id}/members/{usr_id}", headers = {
                    "authorization": f"Bot {login.tkn_value}"
                })
            except:
                messagebox.showerror("Error", "Check your internet connection")
            if (response.status_code != 204):
                messagebox.showerror("Error", response.json()["message"])
            elif (response.status_code == 204):
                run(connection())
                messagebox.showinfo("Success", "The user has been successfully kicked!")
        def ban_usr(guild_id: str, usr_id: str):
            try:
                response = put(f"https://discord.com/api/v10/guilds/{guild_id}/bans/{usr_id}", headers = {
                    "authorization": f"Bot {login.tkn_value}"
                })
            except:
                messagebox.showerror("Error", "Check your internet connection")
            if (response.status_code != 204):
                messagebox.showerror("Error", response.json()["message"])
            elif (response.status_code == 204):
                run(connection())
                messagebox.showinfo("Success", "The user has been successfully banned!")
        def unban_usr(guild_id: str, usr_id: str):
            try:
                response = delete(f"https://discord.com/api/v10/guilds/{guild_id}/bans/{usr_id}", headers = {
                    "authorization": f"Bot {login.tkn_value}"
                })
            except:
                messagebox.showerror("Error", "Check your internet connection")
            if (response.status_code != 204):
                messagebox.showerror("Error", response.json()["message"])
            elif (response.status_code == 204):
                run(connection())
                messagebox.showinfo("Success", "The user has been successfully unbanned!")
        if (choice == "Write a message"):
            post_msg(msg_textbox.get("0.0", "end"))
        elif (choice == "Edit a message"):
            edit_msg(msg_id_textbox.get(), msg_textbox.get("0.0", "end"))
        elif (choice == "Delete a message"):
            delete_msg(msg_id_textbox.get())
        elif (choice == "Pin a message"):
            pin_msg(msg_id_textbox.get())
        elif (choice == "Unpin a message"):
            unpin_msg(msg_id_textbox.get())
        elif (choice == "Create a channel"):
            create_chn(guild_id_textbox.get(), chn_name_textbox.get(), chn_type.get())
        elif (choice == "Edit a channel"):
            edit_chn(chn_id_textbox.get(), chn_name_textbox.get())
        elif (choice == "Delete a channel"):
            delete_chn(chn_id_textbox.get())
        elif (choice == "Create a thread"):
            create_thread(msg_id_textbox.get(), thread_name_textbox.get())
        elif (choice == "Kick a user"):
            kick_usr(guild_id_textbox.get(), usr_id_textbox.get())
        elif (choice == "Ban a user"):
            ban_usr(guild_id_textbox.get(), usr_id_textbox.get())
        elif (choice == "Unban a user"):
            unban_usr(guild_id_textbox.get(), usr_id_textbox.get())
    def check_for_updates():
        if (is_alpha == True) :
            messagebox.showerror("Error", "Not allowed with alpha version")
            return
        try:
            github_version = (get("https://api.github.com/repos/lorypelli/http_requests/releases/latest")).json()["tag_name"]
        except:
            messagebox.showerror("Error", "Check your internet connection")
        if (github_version != user_version):
            response = messagebox.askyesno("Update", "A new version is avaible, please update!")
            if (response == True):
                open("https://github.com/LoryPelli/http_requests/releases/latest")
        else:
            messagebox.showinfo("Update", "No updates avaible!")
    CTkLabel(app, text=login.id, font=("Arial", 16)).place(relx=0.01, rely=0)
    username = CTkLabel(app, text=login.username, font=("Arial", 16))
    username.place(relx=0.6, rely=0)
    CTkButton(app, text="Logout", font=("Arial", 16), width=25, height=15, command=logout).place(relx=0.85, rely=0.01)
    CTkButton(app, text="Check for\nUpdates", font=("Arial", 16), width=25, height=15, command=check_for_updates).place(relx=0.82, rely=0.8)
    chn_id_label = CTkLabel(app, text="Insert channel id", font=("Arial", 16))
    chn_id_label.place(relx=0.01, rely=0.1)
    chn_id_textbox = CTkEntry(app, width=250, height=25, font=("Arial", 16))
    chn_id_textbox.place(relx=0.3, rely=0.1)
    guild_id_label = CTkLabel(app, text="Insert guild id", font=("Arial", 16))
    guild_id_textbox = CTkEntry(app, width=250, height=25, font=("Arial", 16))
    CTkLabel(app, text="Select action", font=("Arial", 16)).place(relx=0.01, rely=0.22)
    combobox = CTkComboBox(app, values=["Write a message", "Edit a message", "Pin a message", "Create a channel", "Edit a channel", "Create a thread", "Delete a channel", "Delete a message", "Unpin a message", "Kick a user", "Ban a user", "Unban a user"], state="readonly", variable=StringVar(value="Write a message"), width=250, font=("Arial", 16), dropdown_font=("Arial", 16), justify="center", command=combochoice)
    combobox.place(relx=0.3, rely=0.22)
    chn_type_label = CTkLabel(app, text="Select channel type", font=("Arial", 16))
    chn_type = CTkComboBox(app, values=["Text", "Voice", "Stage", "Announcement", "Forum"], state="readonly", variable=StringVar(value="Text"), width=250, font=("Arial", 16), dropdown_font=("Arial", 16), justify="center", command=confirm)
    msg_label = CTkLabel(app, text="Insert message", font=("Arial", 16))
    msg_label.place(relx=0.01, rely=0.35)
    msg_textbox = CTkTextbox(app, width=250, height=100, font=("Arial", 16), border_width=2)
    msg_textbox.place(relx=0.3, rely=0.35)
    msg_id_label = CTkLabel(app, text="Insert message id", font=("Arial", 16))
    msg_id_textbox = CTkEntry(app, width=250, height=25, font=("Arial", 16))
    confirm_action = CTkButton(app, text="Send", font=("Arial", 16), command=confirm)
    confirm_action.place(relx=0.4, rely=0.8)
    chn_name_label = CTkLabel(app, text="Insert channel name", font=("Arial", 16))
    chn_name_textbox = CTkEntry(app, width=250, height=25, font=("Arial", 16))
    thread_name_label = CTkLabel(app, text="Insert thread name", font=("Arial", 16))
    thread_name_textbox = CTkEntry(app, width=250, height=25, font=("Arial", 16))
    usr_id_label = CTkLabel(app, text="Insert user id", font=("Arial", 16))
    usr_id_textbox = CTkEntry(app, width=250, height=25, font=("Arial", 16))
    return app
if __name__ == "__main__":
    login().mainloop()