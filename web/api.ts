const callElevator = async (floor: number): Promise<void> => {
  fetch(`/call/${floor}`, {
    method: "POST",
  }).catch((err) => {
    console.error("Failed to call elevator", err)
  })
}

const getFloorCount = (): number => 10

const getElevators = async (): Promise<ElevatorType[]> => {
  return fetch("/elevators")
    .then((res) => res.json())
    .catch((err) => {
      console.error("Failed to get elevators", err)
    })
}

export { callElevator, getFloorCount, getElevators }
