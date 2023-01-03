import requests #importo la libreria request che uso per la richiesta POST
import shutil #importo la libreia shutil che mi servirà per rimuovere una cartella
import os #importo la libreria os per controllare se esiste la cartella
import config #importo il file config.py (puoi fare semplicemente invio nell'input e userà i valori del file di configurazione)
if (os.path.exists("./__pycache__")): #se esiste
    shutil.rmtree("./__pycache__") #rimuovo la cartella __pycache__ perchè è inutile
bot_token = input("Insert the bot token ") or config.bot_token  #ricevo il token del bot in input
response = requests.get("https://discord.com/api/auth/login", headers = { #controllo se il token ricevuto in input è valido
    "authorization": "Bot " + bot_token
})
while (response.status_code != 200): #fino a quando è valido lo continuo a chiedere e controllo se quello inserito è valido
    print("There was an error, the bot token isn't valid. Try again!")
    bot_token = input("Insert the bot token ") or config.bot_token
    response = requests.get("https://discord.com/api/auth/login", headers = {
    "authorization": "Bot " + bot_token
})
channel_id = input("Insert the channel id ") or config.channel_id #ricevo l'id del canale in input
response = requests.get("https://discord.com/api/channels/" + channel_id, headers = { #controllo se l'id ricevuto in input è valido
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