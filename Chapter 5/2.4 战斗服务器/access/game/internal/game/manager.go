package game

var (
	gamePool         = make(map[uint32]*Game, 0)
	newGameId uint32 = 10100
)

func AddGames(count int) {
	for i := 0; i < count; i++ {
		gamePool[newGameId] = NewGame(newGameId)
		newGameId += 10
	}
}

func GetGame() *Game {
	// 获取一个 gamePool 中，处于 Idle 状态的 Game
	for _, game := range gamePool {
		if game.Status() == Idle {
			return game
		}
	}
	return nil
}
