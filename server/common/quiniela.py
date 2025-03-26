import logging

from common.utils import store_bets, Bet


class Quiniela:
    def __init__(self):
        self.bets_amount = 0
        pass

    @staticmethod
    def register_bets(msg):
        bets = []
        try:
            for line in msg.split('\n'):
                logging.info(f'action: split_lines | result: success | line: {line}')
                fields = line.split(',')
                logging.info(
                    f'action: apuesta_leida | result: success | apuesta:{fields[0]} {fields[1]} {fields[2]} {fields[3]} {fields[4]} {fields[5]}')
                bets.append(Bet(fields[0], fields[1], fields[2], fields[3], fields[4], fields[5]))
            store_bets(bets)
            logging.info(f'action: apuesta_recibida | result: success | cantidad: {len(bets)}')
        except IndexError:
            logging.info(f'action: apuesta_recibida | result: fail | cantidad: {len(bets)}')