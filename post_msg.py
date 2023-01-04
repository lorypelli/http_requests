from requests import get, post #the library requests it's used to make POST and GET requests
from shutil import rmtree #I use this to remvoe a folder
from os import path #I use this to check if the folder exists
from config import bot_token, channel_id #Values will be used if not null
if (path.exists("./__pycache__")): #if exists
    rmtree("./__pycache__") #delete the __pycache__ folder because it's unuseful for my project
if (bot_token != None): #if the token in the config.py file is not null
    bot_token = bot_token #it will be the token
else :
    bot_token = input("Insert the bot token ") #else the program will receive it from input
response = get("https://discord.com/api/auth/login", headers = { #check if the token is valid
    "authorization": "Bot " + bot_token
})
while (response.status_code != 200): #until is valid it asks for a new token and check again
    print("There was an error, the bot token isn't valid. Try again!")
    bot_token = input("Insert the bot token ")
    response = get("https://discord.com/api/auth/login", headers = {
    "authorization": "Bot " + bot_token
})
if (channel_id != None): #if the id in the config.py file is not null
    channel_id = channel_id #it will be the id
else :
    channel_id = input("Insert the channel id ") #else the program will receive it from input
response = get("https://discord.com/api/channels/" + channel_id, headers = { #check if the id is valid
    "authorization": "Bot " + bot_token
})
while (response.status_code != 200): #until is valid it asks for a new id and check again
    print("There was an error, the channel id isn't valid. Try again!")
    channel_id = input("Insert the channel id ")
    response = get("https://discord.com/api/channels/" + channel_id, headers = {
    "authorization": "Bot " + bot_token
})
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