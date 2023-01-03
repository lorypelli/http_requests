import requests #importo la libreria request che uso per la richiesta POST
import shutil #importo la libreia shutil che mi servirà per rimuovere una cartella
import os #importo la libreria os per controllare se esiste la cartella
import config #importo il file config.py (puoi fare semplicemente invio nell'input e userà i valori del file di configurazione)
if (os.path.exists("./__pycache__")): #se esiste
    shutil.rmtree("./__pycache__") #rimuovo la cartella __pycache__ perchè è inutile
if (config.bot_token != None): #se il token nel file config.py non è nullo
    bot_token = config.bot_token #allora quello sarà il token
else :
    bot_token = input("Insert the bot token ") #altrimenti lo ricevo in input
response = requests.get("https://discord.com/api/auth/login", headers = { #controllo se il token è valido
    "authorization": "Bot " + bot_token
})
while (response.status_code != 200): #fino a quando è valido lo continuo a chiedere e controllo se quello inserito è valido
    print("There was an error, the bot token isn't valid. Try again!")
    bot_token = input("Insert the bot token ")
    response = requests.get("https://discord.com/api/auth/login", headers = {
    "authorization": "Bot " + bot_token
})
if (config.channel_id != None): #se l'id nel file config.py non è nullo
    channel_id = config.channel_id #allora quello sarà l'id
else :
    channel_id = input("Insert the channel id ") #altrimenti lo ricevo in input
response = requests.get("https://discord.com/api/channels/" + channel_id, headers = { #controllo se l'id è valido
    "authorization": "Bot " + bot_token
})
while (response.status_code != 200): #fino a quando è valido lo continuo a chiedere e controllo se quello inserito è valido
    print("There was an error, the channel id isn't valid. Try again!")
    channel_id = input("Insert the channel id ") or config.channel_id
    response = requests.get("https://discord.com/api/channels/" + channel_id, headers = {
    "authorization": "Bot " + bot_token
})
message = input("Insert the message ") #chiedo di inserire un messaggio all'utente
def post(msg: str): #faccio una funzione per semplificare il tutto
    response = requests.post("https://discord.com/api/channels/" + channel_id + "/messages", headers = { #faccio una richiesta con i vari parametri ricevuti in input
        "authorization": "Bot " + bot_token
    }, data = {
        "content": msg
    })
    if (response.status_code != 200): #se fallisce
        print("There was an error, try again!") #lo scrive in output
    else :
        print("The message has been sent successfully!") #altrimenti scrive in outuput che ha inviato il messaggio con successo
post(message) #esegue la funzione dichiarata prima con il parametro richiesto