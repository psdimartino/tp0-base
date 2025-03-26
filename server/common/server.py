import errno
import socket
import logging
import signal

from common.utils import Bet, store_bets
from common.quiniela import Quiniela


class Server:
    def __init__(self, port, listen_backlog):
        # Set SIGTERM handler
        signal.signal(signal.SIGTERM, self.handle_sigterm)
        # Initialize server socket
        self.running = True
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self.quiniela = Quiniela()

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        while self.running:
            try:
                client_sock = self.__accept_new_connection()
                self.__handle_client_connection(client_sock)
            except socket.error as e:
                if e.errno == errno.EBADF:
                    logging.error(f'action: accept_new_connection | result: success | detail: closed connection')
                else:
                    logging.error(f'action: accept_new_connection | result: error | error: {e}')
                return

        self.__close_server_socket()

    @staticmethod
    def __handle_client_connection(client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        bet_amount = 0
        try:
            while True:
                msg = Server.recv(client_sock)
                logging.info(f'action: mensaje_recibido | result: success | msg: {msg}')
                if msg == "end":
                    break
                Quiniela.register_bets(msg)
                Server.send(client_sock, "ok\n")

        except OSError as e:
            logging.info(f'action: mensaje_recibido | result: fail | error: {e}')
        finally:
            client_sock.close()

    @staticmethod
    def recv(client_sock):
        data = b""
        while True:
            chunk = client_sock.recv(80000)
            if not chunk:
                break
            data += chunk
            if b"\n\n" in chunk:
                logging.info(f'action: message_end_received | result: success')
                break
        return data.rstrip().decode()  # Remove trailing newlines/spaces

    @staticmethod
    def send(client_sock, msg):
        data = (msg + "\n").encode('utf-8')
        total_sent = 0
        while total_sent < len(data):
            sent = client_sock.send(data[total_sent:])
            if sent == 0:
                raise RuntimeError("Socket connection broken")
            total_sent += sent

    def __close_server_socket(self):
        if self._server_socket is None:
            return
        self._server_socket.shutdown(socket.SHUT_RDWR)
        self._server_socket.close()
        self._server_socket = None

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c

    def handle_sigterm(self, signum, frame):
        logging.info('Shutting server gracefully with SIGTERM')
        self.running = False
        self.__close_server_socket()
