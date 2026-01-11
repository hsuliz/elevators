import { callElevator } from "./api"

const renderTable = (floorCount: number, elevators: ElevatorType[]): void => {
  const container = document.getElementById("matrix")
  if (!container) {
    throw new Error("Matrix container not found")
  }

  const table = document.createElement("table")

  // header
  const headerRow = document.createElement("tr")

  // floor header
  const elevatorHeader = document.createElement("th")
  elevatorHeader.innerText = "Floor"
  headerRow.appendChild(elevatorHeader)

  // elevators header
  for (let i = 0; i < elevators.length; i++) {
    const tableHeader = document.createElement("th")
    const elevatorId = i + 1
    tableHeader.innerText = "Elevator " + elevatorId
    headerRow.appendChild(tableHeader)
  }
  table.appendChild(headerRow)

  // floor data
  for (let i = floorCount; i >= 0; i--) {
    const rows = document.createElement("tr")

    const floor = document.createElement("td")
    floor.innerText = i.toString()
    floor.dataset.floorBtn = i.toString() // Add data attribute for easy selection
    floor.title = `Call elevator to floor ${i}`
    floor.addEventListener("click", async () => {
      await callElevator(i)
      floor.classList.add("floor-called") // Visual feedback on click
    })

    rows.appendChild(floor)

    for (const elevator of elevators) {
      const elevatorCell = document.createElement("td")
      elevatorCell.dataset.elevatorId = `${elevator.id}`
      elevatorCell.dataset.floor = i.toString()
      rows.appendChild(elevatorCell)
    }

    table.appendChild(rows)
  }

  container.appendChild(table)
}

// id to floor
const elevatorPositions: Record<number, number> = {}

const populateTable = (elevators: ElevatorType[]): void => {
  for (const elevator of elevators) {
    const cell = document.querySelector<HTMLTableCellElement>(
      `td[data-elevator-id="${elevator.id}"][data-floor="${elevator.currentFloor}"]`,
    )

    if (!cell) continue
    cell.classList.add("elevator")
    cell.textContent = "●"
    elevatorPositions[elevator.id] = elevator.currentFloor
  }
}

const updateElevatorCell = (elevator: ElevatorType) => {
  console.log(elevator)

  const floorBtn = document.querySelector<HTMLTableCellElement>(
    `td[data-floor-btn="${elevator.currentFloor}"]`,
  )
  if (floorBtn) {
    floorBtn.classList.remove("floor-called")
  }

  const currentFloor = elevatorPositions[elevator.id]
  const currentCell = document.querySelector<HTMLTableCellElement>(
    `td[data-elevator-id="${elevator.id}"][data-floor="${currentFloor}"]`,
  )
  if (currentCell) {
    currentCell.classList.remove("elevator")
    currentCell.textContent = ""
  }

  const newCell = document.querySelector<HTMLTableCellElement>(
    `td[data-elevator-id="${elevator.id}"][data-floor="${elevator.currentFloor}"]`,
  )
  if (newCell) {
    newCell.classList.add("elevator")
    newCell.textContent = "●"
  }

  elevatorPositions[elevator.id] = elevator.currentFloor
}

export { renderTable, populateTable, updateElevatorCell }
