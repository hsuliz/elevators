const callElevator = async (floor: number): Promise<void> => {
  fetch(`/call/${floor}`, {
    method: "POST",
  }).catch((err) => {
    console.error("Failed to call elevator", err)
  })
}

const getFloors = (): Promise<Map<number, boolean>> => {
  return fetch("/floors")
    .then((res) => res.json())
    .then((json: any) => new Map(Object.entries(json).map(([k, v]) => [Number(k), Boolean(v)])))
    .catch((err) => {
      console.error("Failed to get floor count", err)
    })
}

const getElevators = async (): Promise<ElevatorType[]> => {
  return fetch("/elevators")
    .then((res) => res.json())
    .catch((err) => {
      console.error("Failed to get elevators", err)
    })
}

export { callElevator, getFloors, getElevators }
