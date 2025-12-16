interface ElevatorState {
    ID: number;
    CurrentFloor: number;
    DestinationFloors: number[];
    Status: number;
}

function callElevator(floor: number): void {
    fetch(`http://localhost:8080/call/${floor}`, {
        method: 'POST'
    }).catch(err => {
        console.error('Failed to call elevator', err);
    });
}


function getElevatorCount(): number {
    return 3;
}

function getFloorCount(): number {
    return 10;
}

const matrixDiv = document.getElementById('matrix') as HTMLDivElement;
let ws: WebSocket;

const elevatorStates: Map<number, ElevatorState> = new Map();

function connect(): void {
    ws = new WebSocket('ws://localhost:8080/ws');

    ws.onopen = () => {
        addLog('System: WebSocket connected');
    };

    ws.onmessage = (event: MessageEvent<string>) => {
        addLog(`Server: ${event.data}`);

        try {
            const json = event.data.startsWith('Server:')
                ? event.data.replace('Server:', '').trim()
                : event.data;

            const state: ElevatorState = JSON.parse(json);
            elevatorStates.set(state.ID, state);
            renderMatrix();
        } catch (e) {
            console.error(e);
        }
    };

    ws.onclose = () => {
        addLog('System: WebSocket disconnected, reconnecting...');
        setTimeout(connect, 1000);
    };
}

function renderMatrix(): void {
    const elevators = getElevatorCount();
    const floors = getFloorCount();

    matrixDiv.innerHTML = '';

    matrixDiv.style.gridTemplateColumns =
        `80px repeat(${elevators}, 60px)`;

    for (let floor = floors; floor >= 0; floor--) {
        const label = document.createElement('div');
        label.className = 'cell label';
        label.textContent = `Floor ${floor}`;

        label.addEventListener('click', () => {
            callElevator(floor);
        });

        matrixDiv.appendChild(label);

        for (let id = 1; id <= elevators; id++) {
            const cell = document.createElement('div');
            cell.className = 'cell';

            const state = elevatorStates.get(id);

            if (state?.CurrentFloor === floor) {
                cell.classList.add('elevator');
                cell.textContent = 'E';
            } else if (state?.DestinationFloors.includes(floor)) {
                cell.classList.add('destination');
                cell.textContent = 'D';
            }

            matrixDiv.appendChild(cell);
        }
    }
}

const logsDiv = document.getElementById('logs') as HTMLDivElement;

function addLog(message: string): void {
    const entry = document.createElement('div');
    entry.className = 'log-entry';
    entry.textContent = message;

    logsDiv.appendChild(entry);
    logsDiv.scrollTop = logsDiv.scrollHeight;
}


window.addEventListener('DOMContentLoaded', () => {
    renderMatrix()
    connect()
});
