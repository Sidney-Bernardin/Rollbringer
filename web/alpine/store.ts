const { gameId, name } = document.body.dataset
if (gameId) {
    window.alpine.store('game', {
        id: gameId,
        name: name,
    })
}

window.alpine.store('rolling', false)
