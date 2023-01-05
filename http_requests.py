from requests import get, post, delete #the library requests it's used to make POST, GET and DELETE requests
from pick import pick
from shutil import rmtree #I use this to remove a folder
from os import path #I use this to check if the folder exists
from config import bot_token, channel_id #these values will be used if not null
if (path.exists("./__pycache__")): #if exists
    rmtree("./__pycache__") #delete the __pycache__ folder because it's unuseful for my project
response = get("https://discord.com/api/auth/login", headers = { #check if the token is valid
    "authorization": "Bot " + bot_token
})
if (response.status_code == 200): #if the token in the config.py file is valid
    bot_token = bot_token #it will be the token
    print("The bot token has been taken from the config.py file")
else :
    bot_token = input("Insert the bot token ") #else the program will receive it from input
response = get("https://discord.com/api/auth/login", headers = { #check if the token is valid
    "authorization": "Bot " + bot_token
})
while (response.status_code != 200): #until is valid it asks for a new token and check again
    print("There was an error, try again!")
    bot_token = input("Insert the bot token ")
    response = get("https://discord.com/api/auth/login", headers = {
    "authorization": "Bot " + bot_token
})
response = get("https://discord.com/api/channels/" + channel_id, headers = {
    "authorization": "Bot " + bot_token
})
if (response.status_code == 200): #if the id in the config.py file is valid
    channel_id = channel_id #it will be the id
    print("The channel id has been taken from the config.py file")
else :
    channel_id = input("Insert the channel id ") #else the program will receive it from input
response = get("https://discord.com/api/channels/" + channel_id, headers = { #check if the id is valid
    "authorization": "Bot " + bot_token
})
while (response.status_code != 200): #until is valid it asks for a new id and check again
    print("There was an error, try again!")
    channel_id = input("Insert the channel id ")
    response = get("https://discord.com/api/channels/" + channel_id, headers = {
    "authorization": "Bot " + bot_token
})
menu_title = "Do you want to write a new message or delete a channel? (J / K to move)" #the menu title
menu_options = ["Write a message", "Delete a channel (this will use channel id you already provided)"] #the menu options
option, index = pick(menu_options, menu_title, indicator="–>") #pick an option
print(option) #write the option
if (index == 1): #if option index it's 1 (not 2 because an array starts with 0)
    def delete_chn(chn_id: str): #use a function because is better
        response = delete("https://discord.com/api/channels/" + chn_id, headers = { #make a request with input params
            "authorization": "Bot " + bot_token
        })
        if (response.status_code != 200): #if it fails
            print(response.status_code)
            print("There was an error, try again!") #write that there was an error
        else :
            print("The channel has been deleted successfully!") #else write that the channel has been deleted successfully
    delete_chn(channel_id) #run the function with the requested param
else: #else option it's equal to 0
    message = input("Insert the message ") #asks to the user for a message input
    def post_msg(msg: str): #use a function because is better
        response = post("https://discord.com/api/channels/" + channel_id + "/messages", headers = { #make a request with input params
            "authorization": "Bot " + bot_token
        }, data = {
            "content": msg
        })
        if (response.status_code != 200): #if it fails
            print("There was an error, try again!") #write that there was an error
        else :
            print("The message has been sent successfully!") #else write that the message has been sent successfully
    post_msg(message) #run the function with the requested param