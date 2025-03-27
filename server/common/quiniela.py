import logging
import threading

from common.utils import store_bets, Bet, load_bets, has_won


class Quiniela:
    def __init__(self):
        self.bets_amount = 0
        self.lock = threading.Lock()
        self.bets = None
        self.event = threading.Event()
        pass

    def register_bets(self, msg):
        bets = []
        try:
            for line in msg.split('\n'):
                logging.info(f'action: split_lines | result: success | line: {line}')
                fields = line.split(',')
                logging.info(
                    f'action: apuesta_leida | result: success | apuesta:{fields[0]} {fields[1]} {fields[2]} {fields[3]} {fields[4]} {fields[5]}')
                bets.append(Bet(fields[0], fields[1], fields[2], fields[3], fields[4], fields[5]))
            with self.lock:
                store_bets(bets)
            logging.info(f'action: apuesta_recibida | result: success | cantidad: {len(bets)}')
        except IndexError:
            logging.info(f'action: apuesta_recibida | result: fail | cantidad: {len(bets)}')

    def load_bets(self):
        with self.lock:
            if self.bets is None:
                logging.info(f'action: carga_apuestas | result: in_progress')
                self.bets = list(load_bets())
                logging.info(f'action: carga_apuestas | result: success')
                self.event.set()
            logging.info(f'action: obtener_apuestas | result: success | cantidad: {len(self.bets)}')
            return self.bets

    def winners(self, agency):
        self.event.wait()
        logging.info(f'action: consulta_ganadores | result: in_progress | agencia: agency')
        return [b for b in self.bets if b.agency == agency and has_won(b)]
