let elevators: ElevatorType[] = [];

const connect = (): void => {
    const ws = new WebSocket('ws://localhost:8080/ws');

    ws.onopen = () => {
        console.log("WebSocket: connected");
    };

    ws.onmessage = (event: MessageEvent<string>) => {
        try {
            elevators = JSON.parse(event.data);
            console.log("Got elevators:", elevators)

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
