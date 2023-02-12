import customtkinter
from tkinter import messagebox
from requests import get, post, delete, patch, put
from webbrowser import open
from urllib import request
from os import environ, getenv, path, makedirs
user_version = "ThirteenthRelease"
is_alpha = False
icon_url = "https://raw.githubusercontent.com/LoryPelli/http_requests/main/app_icon.ico"
icon_directory = f"{environ.get('SystemDrive')}/Users/{getenv('Username')}/http_requests"
icon_file = f"{icon_directory}/app_icon.ico"
customtkinter.set_appearance_mode("System")
customtkinter.set_default_color_theme("dark-blue")
if (not(path.exists(icon_directory))):
    makedirs(icon_directory)
try:
    request.urlretrieve(icon_url, icon_file)
except:
    pass
def login():
    app = customtkinter.CTk()
    app.title("Login")
    app.resizable(False, False)
    app.geometry("840x70+500+500")
    try:
        app.wm_iconbitmap(icon_file)
    except:
        pass
    def loginbtn():
        try:
            response = get("https://discord.com/api/auth/login", headers = {
                "authorization": f"Bot {tkn_textbox.get()}"
            })
        except:
            messagebox.showerror("Error", "Check your internet connection")
        if (response.status_code != 200):
            messagebox.showerror("Error", "There was an error, try again!")
        elif (response.status_code == 200):
            messagebox.showinfo("Success", "Validation Passed!")
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
    try:
        app.wm_iconbitmap(icon_file)
    except:
        pass
    def logout():
        response = messagebox.askyesno("Logout", "Are you sure you want to logout")
        if (response == True):
            app.destroy()
            login().mainloop()
    def combochoice(choice: str):
        if (choice == "Write a message"):
            msg_label.place(relx=0.01, rely=0.35)
            msg_textbox.place(relx=0.3, rely=0.35)
            chn_name_label.place_forget()
            chn_name_textbox.place_forget()
            msg_id_label.place_forget()
            msg_id_textbox.place_forget()
            thread_name_label.place_forget()
            thread_name_textbox.place_forget()
            confirm_action.configure(text="Send")
        elif (choice == "Edit a channel"):
            chn_name_label.place(relx=0.01, rely=0.35)
            chn_name_textbox.place(relx=0.3, rely=0.35)
            msg_label.place_forget()
            msg_textbox.place_forget()
            msg_id_label.place_forget()
            msg_id_textbox.place_forget()
            thread_name_label.place_forget()
            thread_name_textbox.place_forget()
            confirm_action.configure(text="Edit")
        elif (choice == "Edit a message"):
            msg_id_label.place(relx=0.01, rely=0.35)
            msg_id_textbox.place(relx=0.3, rely=0.35)
            msg_label.place(relx=0.01, rely=0.47)
            msg_textbox.place(relx=0.3, rely=0.47)
            msg_textbox.configure(height=75)
            chn_name_label.place_forget()
            chn_name_textbox.place_forget()
            thread_name_label.place_forget()
            thread_name_textbox.place_forget()
            confirm_action.configure(text="Edit")
        elif (choice == "Pin a message"):
            msg_id_label.place(relx=0.01, rely=0.35)
            msg_id_textbox.place(relx=0.3, rely=0.35)
            msg_label.place_forget()
            msg_textbox.place_forget()
            chn_name_label.place_forget()
            chn_name_textbox.place_forget()
            thread_name_label.place_forget()
            thread_name_textbox.place_forget()
            confirm_action.configure(text="Pin")
        elif (choice == "Unpin a message"):
            msg_id_label.place(relx=0.01, rely=0.35)
            msg_id_textbox.place(relx=0.3, rely=0.35)
            msg_label.place_forget()
            msg_textbox.place_forget()
            chn_name_label.place_forget()
            chn_name_textbox.place_forget()
            thread_name_label.place_forget()
            thread_name_textbox.place_forget()
            confirm_action.configure(text="Unpin")
        elif (choice == "Delete a message"):
            msg_id_label.place(relx=0.01, rely=0.35)
            msg_id_textbox.place(relx=0.3, rely=0.35)
            msg_label.place_forget()
            msg_textbox.place_forget()
            chn_name_label.place_forget()
            chn_name_textbox.place_forget()
            thread_name_label.place_forget()
            thread_name_textbox.place_forget()
            confirm_action.configure(text="Delete")
        elif (choice == "Delete a channel"):
            msg_id_label.place_forget()
            msg_id_textbox.place_forget()
            msg_label.place_forget()
            msg_textbox.place_forget()
            chn_name_label.place_forget()
            chn_name_textbox.place_forget()
            thread_name_label.place_forget()
            thread_name_textbox.place_forget()
            confirm_action.configure(text="Delete")
        elif (choice == "Create a thread"):
            msg_id_label.place(relx=0.01, rely=0.35)
            msg_id_textbox.place(relx=0.3, rely=0.35)
            thread_name_label.place(relx=0.01, rely=0.465)
            thread_name_textbox.place(relx=0.3, rely=0.465)
            msg_label.place_forget()
            msg_textbox.place_forget()
            chn_name_label.place_forget()
            chn_name_textbox.place_forget()
            confirm_action.configure(text="Create")
    def confirm():
        choice = combobox.get()
        def post_msg(msg: str):
            try:
                response = post(f"https://discord.com/api/channels/{chn_id_textbox.get()}/messages", headers = {
                    "authorization": f"Bot {login.tkn_value}"
                }, json = {
                    "content": msg
                })
            except:
                messagebox.showerror("Error", "Check your internet connection")
            if (response.status_code != 200):
                messagebox.showerror("Error", "There was an error, try again!")
            elif (response.status_code == 200):
                messagebox.showinfo("Success", "The message has been sent successfully!")
        def edit_msg(msg_id: str, msg: str):
            response = patch(f"https://discord.com/api/channels/{chn_id_textbox.get()}/messages/{msg_id}", headers = {
                "authorization": f"Bot {login.tkn_value}"
            }, json = {
                "content": msg
            })
            if (response.status_code != 200):
                messagebox.showerror("Error", "There was an error, try again!")
            elif (response.status_code == 200):
                messagebox.showinfo("Success", "The message has been edited successfully!")
        def delete_msg(msg_id: str):
            response = delete(f"https://discord.com/api/channels/{chn_id_textbox.get()}/messages/{msg_id}", headers = {
                "authorization": f"Bot {login.tkn_value}"
            })
            if (response.status_code != 204):
                messagebox.showerror("Error", "There was an error, try again!")
            elif (response.status_code == 204):
                messagebox.showinfo("Success", "The message has been deleted successfully!")
        def pin_msg(msg_id: str):
            try:
                response = put(f"https://discord.com/api/channels/{chn_id_textbox.get()}/pins/{msg_id}", headers = {
                    "authorization": f"Bot {login.tkn_value}"
                })
            except:
                messagebox.showerror("Error", "Check your internet connection")
            if (response.status_code != 204):
                messagebox.showerror("Error", "There was an error, try again!")
            elif (response.status_code == 204):
                messagebox.showinfo("Success", "The message has been pinned successfully!")
        def unpin_msg(msg_id: str):
            try:
                response = delete(f"https://discord.com/api/channels/{chn_id_textbox.get()}/pins/{msg_id}", headers = {
                    "authorization": f"Bot {login.tkn_value}"
                })
            except:
                messagebox.showerror("Error", "Check your internet connection")
            if (response.status_code != 204):
                messagebox.showerror("Error", "There was an error, try again!")
            elif (response.status_code == 204):
                messagebox.showinfo("Success", "The message has been unpinned successfully!")
        def edit_chn(chn_id: str, chn_name: str):
            try:
                response = patch(f"https://discord.com/api/channels/{chn_id}", headers = {
                    "authorization": f"Bot {login.tkn_value}"
                }, json = {
                    "name": chn_name
                })
            except:
                messagebox.showerror("Error", "Check your internet connection")
            if (response.status_code != 200):
                messagebox.showerror("Error", "There was an error, try again!")
            elif (response.status_code == 200):
                messagebox.showinfo("Success", "The channel has been edited successfully!")
        def delete_chn(chn_id: str):
            try:
                response = delete(f"https://discord.com/api/channels/{chn_id}", headers = {
                    "authorization": f"Bot {login.tkn_value}"
                })
            except:
                messagebox.showerror("Error", "Check your internet connection")
            if (response.status_code != 200):
                messagebox.showerror("Error", "There was an error, try again!")
            elif (response.status_code == 200):
                messagebox.showinfo("Success", "The channel has been deleted successfully!")
        def create_thread(msg_id: str, thread_name: str):
            try:
                response = post(f"https://discord.com/api/channels/{chn_id_textbox.get()}/messages/{msg_id}/threads", headers = {
                    "authorization": "Bot " + login.tkn_value
                }, json = {
                    "name": thread_name
                })
            except:
                messagebox.showerror("Error", "Check your internet connection")
            if (response.status_code != 201):
                messagebox.showerror("Error", "There was an error, try again!")
            elif (response.status_code == 201):
                messagebox.showinfo("Success", "The thread has been created successfully!")
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
        elif (choice == "Edit a channel"):
            edit_chn(chn_id_textbox.get(), chn_name_textbox.get())
        elif (choice == "Delete a channel"):
            delete_chn(chn_id_textbox.get())
        elif (choice == "Create a thread"):
            create_thread(msg_id_textbox.get(), thread_name_textbox.get())
    def check_for_updates():
        try:
            github_version = (get("https://api.github.com/repos/lorypelli/http_requests/releases/latest")).json()["tag_name"]
        except:
            messagebox.showerror("Error", "Check your internet connection")
        if (github_version != user_version and is_alpha == False):
            response = messagebox.askyesno("Update", "A new version is avaible, please update!")
            if (response == True):
                open("https://github.com/LoryPelli/http_requests/releases/latest")
        else:
            messagebox.showinfo("Update", "No updates avaible!")
    customtkinter.CTkLabel(app, text=login.id, font=("Arial", 16)).place(relx=0.01, rely=0)
    username = customtkinter.CTkLabel(app, text=login.username, font=("Arial", 16))
    username.place(relx=0.6, rely=0)
    customtkinter.CTkButton(app, text="Logout", font=("Arial", 16), width=25, height=15, command=logout).place(relx=0.85, rely=0.01)
    customtkinter.CTkButton(app, text="Check for\nUpdates", font=("Arial", 16), width=25, height=15, command=check_for_updates).place(relx=0.82, rely=0.8)
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
    msg_id_label = customtkinter.CTkLabel(app, text="Insert message id", font=("Arial", 16))
    msg_id_textbox = customtkinter.CTkEntry(app, width=250, height=25, font=("Arial", 16))
    confirm_action = customtkinter.CTkButton(app, text="Send", font=("Arial", 16), command=confirm)
    confirm_action.place(relx=0.4, rely=0.8)
    chn_name_label = customtkinter.CTkLabel(app, text="Insert channel name", font=("Arial", 16))
    chn_name_textbox = customtkinter.CTkEntry(app, width=250, height=25, font=("Arial", 16))
    thread_name_label = customtkinter.CTkLabel(app, text="Insert thread name", font=("Arial", 16))
    thread_name_textbox = customtkinter.CTkEntry(app, width=250, height=25, font=("Arial", 16))
    return app
if __name__ == "__main__":
    login().mainloop()