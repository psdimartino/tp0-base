import errno
import socket
import logging
import signal

from common.utils import store_bets, Bet


class Quiniela:
    def __init__(self):
        pass

    @staticmethod
    def register_bet(agency: str, first_name: str, last_name: str, document: str, birthdate: str, number: str):
        bet = Bet(agency, first_name, last_name, document, birthdate, number)
        store_bets([bet])
        logging.info(f'action: apuesta_almacenada | result: success | dni: {document} | numero: {number}')
