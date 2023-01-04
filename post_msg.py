from requests import get, post #importo la libreria request che uso per la richiesta POST
from shutil import rmtree #importo la libreia shutil che mi servirà per rimuovere una cartella
from os import path #importo la libreria os per controllare se esiste la cartella
from config import bot_token, channel_id #importo il file config.py (userà i valori del file di configurazione se non nulli, altrimenti li chiederà all'utente)
if (path.exists("./__pycache__")): #se esiste
    rmtree("./__pycache__") #rimuovo la cartella __pycache__ perchè è inutile per il mio progetto, se fosse un progetto grande, la terrei
if (bot_token != None): #se il token nel file config.py non è nullo
    bot_token = bot_token #allora quello sarà il token
else :
    bot_token = input("Insert the bot token ") #altrimenti lo ricevo in input
response = get("https://discord.com/api/auth/login", headers = { #controllo se il token è valido
    "authorization": "Bot " + bot_token
})
while (response.status_code != 200): #fino a quando è valido lo continuo a chiedere e controllo se quello inserito è valido
    print("There was an error, the bot token isn't valid. Try again!")
    bot_token = input("Insert the bot token ")
    response = get("https://discord.com/api/auth/login", headers = {
    "authorization": "Bot " + bot_token
})
if (channel_id != None): #se l'id nel file config.py non è nullo
    channel_id = channel_id #allora quello sarà l'id
else :
    channel_id = input("Insert the channel id ") #altrimenti lo ricevo in input
response = get("https://discord.com/api/channels/" + channel_id, headers = { #controllo se l'id è valido
    "authorization": "Bot " + bot_token
})
while (response.status_code != 200): #fino a quando è valido lo continuo a chiedere e controllo se quello inserito è valido
    print("There was an error, the channel id isn't valid. Try again!")
    channel_id = input("Insert the channel id ")
    response = get("https://discord.com/api/channels/" + channel_id, headers = {
    "authorization": "Bot " + bot_token
})
message = input("Insert the message ") #chiedo di inserire un messaggio all'utente
def post_msg(msg: str): #faccio una funzione per semplificare il tutto
    response = post("https://discord.com/api/channels/" + channel_id + "/messages", headers = { #faccio una richiesta con i vari parametri ricevuti in input
        "authorization": "Bot " + bot_token
    }, data = {
        "content": msg
    })
    if (response.status_code != 200): #se fallisce
        print("There was an error, try again!") #lo scrive in output
    else :
        print("The message has been sent successfully!") #altrimenti scrive in outuput che ha inviato il messaggio con successo
post_msg(message) #esegue la funzione dichiarata prima con il parametro richiesto