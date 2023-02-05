import customtkinter
import tkinter.messagebox
from requests import get, post, delete, patch, put
from webbrowser import open
user_version = "NinthRelease"
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
    app.geometry("500x250+450+450")
    app.wm_iconbitmap("app_icon.ico")
    def logout():
        response = tkinter.messagebox.askyesno("Logout", "Are you sure you want to logout")
        if (response == True):
            app.destroy()
            login().mainloop()
    def combochoice(choice: str):
        if (choice == "Write a message"):
            msg_label.place(relx=0.01, rely=0.35)
            msg_textbox.place(relx=0.3, rely=0.35)
            chn_name_label.place_forget()
            chn_name_textbox.place_forget()
        elif (choice == "Edit a channel"):
            msg_label.place_forget()
            msg_textbox.place_forget()
            chn_name_label.place(relx=0.01, rely=0.35)
            chn_name_textbox.place(relx=0.3, rely=0.35)
            confirm_action.configure(text="Edit")
        return choice
    def confirm():
        def post_msg(msg: str):
            response = post("https://discord.com/api/channels/" + chn_id_textbox.get() + "/messages", headers = {
                "authorization": "Bot " + login.tkn_value
            }, json = {
                "content": msg
            })
            if (response.status_code != 200):
                tkinter.messagebox.showerror("Error", "There was an error, try again!")
            elif (response.status_code == 200):
                tkinter.messagebox.showinfo("Success", "The message has been sent successfully!")
        def edit_chn(chn_id: str, chn_name: str):
            response = patch("https://discord.com/api/channels/" + chn_id, headers = {
                "authorization": "Bot " + login.tkn_value
            }, json = {
                "name": chn_name
            })
            if (response.status_code != 200):
                tkinter.messagebox.showerror("Error", "There was an error, try again!")
            elif (response.status_code == 200):
                tkinter.messagebox.showinfo("Success", "The channel has been edited successfully!")
        if (combochoice("Write a message")):
            post_msg(msg_textbox.get("0.0", "end"))
        elif (combochoice("Edit a channel")):
            edit_chn(chn_id_textbox, chn_name_textbox)
    customtkinter.CTkLabel(app, text=login.id, font=("Arial", 16)).place(relx=0.01, rely=0)
    username = customtkinter.CTkLabel(app, text=login.username, font=("Arial", 16))
    username.place(relx=0.65, rely=0)
    customtkinter.CTkButton(app, text="Logout", font=("Arial", 16), width=25, height=15, command=logout).place(relx=0.85, rely=0.01)
    customtkinter.CTkLabel(app, text="Insert channel id", font=("Arial", 16)).place(relx=0.01, rely=0.1)
    chn_id_textbox = customtkinter.CTkEntry(app, width=250, height=25, font=("Arial", 16))
    chn_id_textbox.place(relx=0.3, rely=0.1)
    customtkinter.CTkLabel(app, text="Select action", font=("Arial", 16)).place(relx=0.01, rely=0.22)
    combobox = customtkinter.CTkComboBox(app, values=["Write a message", "Edit a message", "Pin a message", "Edit a channel", "Create a thread", "Delete a channel", "Delete a message", "Unpin a message"], state="readonly", variable=customtkinter.StringVar(value="Write a message"), width=250, font=("Arial", 16), dropdown_font=("Arial", 16), justify="center", command=combochoice)
    combobox.place(relx=0.3, rely=0.22)
    msg_label = customtkinter.CTkLabel(app, text="Insert message", font=("Arial", 16))
    msg_label.place(relx=0.01, rely=0.35)
    msg_textbox = customtkinter.CTkTextbox(app, width=250, height=100, font=("Arial", 16), border_width=2)
    msg_textbox.place(relx=0.3, rely=0.35)
    confirm_action = customtkinter.CTkButton(app, text="Send", font=("Arial", 16), command=confirm)
    confirm_action.place(relx=0.4, rely=0.8)
    chn_name_label = customtkinter.CTkLabel(app, text="Insert channel name", font=("Arial", 16))
    chn_name_textbox = customtkinter.CTkEntry(app, width=250, height=25, font=("Arial", 16))
    return app
if __name__ == "__main__":
    if (user_version != github_version and is_alpha == False):
        response = tkinter.messagebox.showinfo("info", "A new version is avaible, please update!")
        if (response == "ok"):
            open("https://github.com/LoryPelli/http_requests/releases/latest")
    else:
        login().mainloop()