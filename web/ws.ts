import {updateElevatorCell} from "./render";

const connect = (): void => {
    const ws = new WebSocket('ws://localhost:8080/ws');

    ws.onopen = () => {
        console.log("WebSocket: connected");
    };

    ws.onmessage = (event: MessageEvent<string>) => {
        try {
            const elevator: ElevatorType = JSON.parse(event.data);
            updateElevatorCell(elevator);
        } catch (e) {
            console.error('Failed to update elevators', e);
        }
    };

    ws.onclose = () => {
        console.log("WebSocket: disconnected, reconnecting...");
        setTimeout(connect, 1000);
    };
};

export {connect};
