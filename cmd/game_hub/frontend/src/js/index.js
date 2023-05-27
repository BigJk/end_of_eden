import 'xterm/css/xterm.css';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';

let CLIENT_ID = '';
let WS_CONN = null;

function getCookie(cname) {
	let name = cname + '=';
	let decodedCookie = decodeURIComponent(document.cookie);
	let ca = decodedCookie.split(';');
	for (let i = 0; i < ca.length; i++) {
		let c = ca[i];
		while (c.charAt(0) == ' ') {
			c = c.substring(1);
		}
		if (c.indexOf(name) == 0) {
			return c.substring(name.length, c.length);
		}
	}
	return '';
}

const term = new Terminal();
const fitAddon = new FitAddon();

term.onData((chunk) => {
	WS_CONN.send(chunk);
});

term.onResize((size) => {
	console.log(size);

	fetch(`/api/resize`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify(size),
	});
});

term.loadAddon(fitAddon);
term.open(document.getElementById('terminal'));
fitAddon.fit();

fetch('/api/create').then((res) => {
	CLIENT_ID = getCookie('id');

	let ws = new WebSocket(`ws://${window.location.host}/ws/${CLIENT_ID}`);

	ws.onopen = () => {
		console.log('Connected to websocket');
	};

	ws.onmessage = (e) => {
		term.write(e.data);
	};

	ws.onclose = () => {
		console.log('Disconnected from websocket');
	};

	WS_CONN = ws;
});
