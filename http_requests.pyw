import customtkinter
import tkinter.messagebox
from requests import get, post, delete, patch, put
from shutil import rmtree
from os import path
from importlib import reload
from psutil import process_iter
from webbrowser import open
user_version = "EighthRelease"
is_alpha = True
github_version = (get("https://api.github.com/repos/lorypelli/http_requests/releases/latest")).json()["tag_name"]
customtkinter.set_appearance_mode("System")
customtkinter.set_default_color_theme("dark-blue")
def login():
    app = customtkinter.CTk()
    app.title("Login")
    app.resizable(False, False)
    app.geometry("840x140+500+500")
    app.wm_iconbitmap("app_icon.ico")
    def config():
            for process in process_iter():
                if (process.name() == "http_requests.exe"):
                    tkinter.messagebox.showerror("Error", "This doesn't work with executable files")
                    break
            else:
                try:
                    import config
                    reload(config)
                    configtkn = config.config()
                    if (path.exists("./__pycache__")):
                        rmtree("./__pycache__")
                    response = get("https://discord.com/api/auth/login", headers = {
                        "authorization": "Bot " + configtkn
                    })
                    if (response.status_code != 200):
                        tkinter.messagebox.showerror("Error", "There was an error, try again!")
                    elif (response.status_code == 200):
                        tkn_textbox.insert(0, configtkn)
                except:
                    tkinter.messagebox.showerror("Error", "The file config.py doesn't exists or is not configured as it should be")
    def loginbtn():
        response = get("https://discord.com/api/auth/login", headers = {
            "authorization": "Bot " + tkn_textbox.get()
        })
        if (response.status_code != 200):
            tkinter.messagebox.showerror("Error", "There was an error, try again!")
        elif (response.status_code == 200):
            tkinter.messagebox.showinfo("Success", "Validation Passed!")
            login.tkn_value = tkn_textbox.get()
            login.username = (get("https://discord.com/api/users/@me", headers = {
                "authorization": "Bot " + tkn_textbox.get()
            })).json()["username"]
            login.id = (get("https://discord.com/api/users/@me", headers = {
                "authorization": "Bot " + tkn_textbox.get()
            })).json()["id"]
            app.destroy()
    customtkinter.CTkButton(app, text="Import from config", command=config, font=("Arial", 16)).place(relx=0.4, rely=0.1)
    customtkinter.CTkLabel(app, text="Insert bot token", font=("Arial", 16)).place(relx=0.01, rely=0.40)
    tkn_textbox = customtkinter.CTkEntry(app, width=700, height=25, font=("Arial", 16))
    tkn_textbox.place(relx=0.15, rely=0.40)
    customtkinter.CTkButton(app, text="Validate", command=loginbtn, font=("Arial", 16)).place(relx=0.4, rely=0.70)
    return app
if __name__ == "__main__":
    if (user_version != github_version and is_alpha == False):
        tkinter.messagebox.showinfo("info", "A new version is avaible, please update!")
        open("https://github.com/LoryPelli/http_requests/releases/latest")
    else:
        login().mainloop()