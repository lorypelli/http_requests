import customtkinter
import tkinter.messagebox
from requests import get, post, delete, patch, put
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
    app.geometry("840x70+500+500")
    app.wm_iconbitmap("app_icon.ico")
    def loginbtn():
        response = get("https://discord.com/api/auth/login", headers = {
            "authorization": f"Bot {tkn_textbox.get()}"
        })
        if (response.status_code != 200):
            tkinter.messagebox.showerror("Error", "There was an error, try again!")
        elif (response.status_code == 200):
            tkinter.messagebox.showinfo("Success", "Validation Passed!")
            login.tkn_value = tkn_textbox.get()
            login.username = (get("https://discord.com/api/users/@me", headers = {
                "authorization": f"Bot {tkn_textbox.get()}"
            })).json()["username"]
            login.id = (get("https://discord.com/api/users/@me", headers = {
                "authorization": f"Bot {tkn_textbox.get()}"
            })).json()["id"]
            app.destroy()
            program().mainloop()
    customtkinter.CTkLabel(app, text="Insert bot token", font=("Arial", 16)).place(relx=0.01, rely=0.1)
    tkn_textbox = customtkinter.CTkEntry(app, width=700, height=25, font=("Arial", 16))
    tkn_textbox.place(relx=0.15, rely=0.1)
    customtkinter.CTkButton(app, text="Validate", command=loginbtn, font=("Arial", 16)).place(relx=0.4, rely=0.5)
    return app
def program():
    app = customtkinter.CTk()
    app.title("http_requests")
    app.resizable(False, False)
    app.geometry("500x420+450+450")
    app.wm_iconbitmap("app_icon.ico")
    def chnidbtn():
        response = get(f"https://discord.com/api/channels/{chn_id_textbox.get()}", headers = {
            "authorization": f"Bot {login.tkn_value}"
        })
        if (response.status_code != 200):
            tkinter.messagebox.showerror("Error", "There was an error, try again!")
        elif (response.status_code == 200):
            tkinter.messagebox.showinfo("Success", "Validation Passed!")
    customtkinter.CTkLabel(app, text=login.id, font=("Arial", 16)).place(relx=0.01, rely=0)
    username = customtkinter.CTkLabel(app, text=login.username, font=("Arial", 16))
    username.place(relx=0.65, rely=0)
    customtkinter.CTkButton(app, text="Logout", font=("Arial", 16)).place(relx=0.36, rely=0.01)
    customtkinter.CTkLabel(app, text="Insert channel id", font=("Arial", 16)).place(relx=0.01, rely=0.1)
    chn_id_textbox = customtkinter.CTkEntry(app, width=250, height=25, font=("Arial", 16))
    chn_id_textbox.place(relx=0.3, rely=0.1)
    customtkinter.CTkButton(app, text="Validate", command=chnidbtn, font=("Arial", 16)).place(relx=0.4, rely=0.18)
    customtkinter.CTkLabel(app, text="Select action", font=("Arial", 16)).place(relx=0.01, rely=0.3)
    combobox = customtkinter.CTkComboBox(app, values=["Write a message", "Edit a message", "Pin a message", "Edit a channel", "Create a thread", "Delete a channel", "Delete a message", "Unpin a message"], state="readonly", variable=customtkinter.StringVar(value="Write a message"))
    combobox.place(relx=0.3, rely=0.3)
    return app
if __name__ == "__main__":
    if (user_version != github_version and is_alpha == False):
        response = tkinter.messagebox.showinfo("info", "A new version is avaible, please update!")
        if (response == "ok"):
            open("https://github.com/LoryPelli/http_requests/releases/latest")
    else:
        login().mainloop()