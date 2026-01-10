import {connect} from "./ws"
import {getElevators, getFloorCount} from "./api";
import {populateTable, renderTable} from "./render";


window.addEventListener('DOMContentLoaded', async () => {
    try {
        const elevators = await getElevators()
        const floorCount = getFloorCount()

        renderTable(floorCount, elevators)
        populateTable(elevators)
        connect()
    } catch (err) {
        console.error(err)
    }
})
