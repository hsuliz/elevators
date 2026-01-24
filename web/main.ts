import { connect } from "./ws"
import { getElevators, getFloors } from "./api"
import { populateTable, renderTable } from "./render"

window.addEventListener("DOMContentLoaded", async () => {
  try {
    const elevators = await getElevators()
    const floors = await getFloors()

    renderTable(floors.size, elevators)
    populateTable(elevators)
    connect()
  } catch (err) {
    console.error(err)
  }
})
