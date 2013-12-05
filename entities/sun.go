package entities

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/Vladimiroff/vec2d"
)

type Sun struct {
	Username string
	Name     string
	speed    int32
	target   *vec2d.Vector
	Position *vec2d.Vector
}

// Database key.
func (s *Sun) Key() string {
	return fmt.Sprintf("sun.%s", s.Name)
}

// Returns the set by X or Y where this entity has to be put in
func (s *Sun) AreaSet() string {
	return fmt.Sprintf(
		ENTITIES_AREA_TEMPLATE,
		RoundCoordinateTo(s.Position.X),
		RoundCoordinateTo(s.Position.Y),
	)
}

// Generate sun's name out of user's initials and 3-digit random number
func (s *Sun) generateName(nickname string) {
	hash, _ := strconv.ParseInt(GenerateHash(nickname)[0:18], 10, 64)
	random := rand.New(rand.NewSource(hash))
	initials := extractUsernameInitials(nickname)
	number := random.Int31n(899) + 100 // we need a 3-digit number
	s.Name = fmt.Sprintf("%s%v", initials, number)
}

// Generates the key of the start position node
func (s *Sun) getStartPosNode(friends *[]Sun) *Node {
	targetPosition := vec2d.New(0, 0)

	//Find best position between all friends
	for _, friend := range friends {
		targetPosition.Collect(friend.Position)
	}
	targetPosition.DivToFloat64(float64(len(friends)))

	//Approximate target to nearest node
	targetPosition.X = SUNS_SOLAR_SYSTEM_RADIUS * math.Floor((2*targetPosition.X/SUNS_SOLAR_SYSTEM_RADIUS)+0.5)
	targetPosition.Y = SUNS_SOLAR_SYSTEM_RADIUS * math.Floor((2*targetPosition.Y/SUNS_SOLAR_SYSTEM_RADIUS*math.Sqrt(3))+0.5)

	return newNode(targetPosition.X, targetPosition.Y)

}

func (s *Sun) fetchNodesLayer(rootNode *Node, zLevel uint32) (results []*Node) {
	horizontalOffset := math.Floor((SUNS_SOLAR_SYSTEM_RADIUS / 2) + 0.5)
	verticalOffset := math.Floor((SUNS_SOLAR_SYSTEM_RADIUS * math.Sqrt(3) / 2) + 0.5)

	results = append(results, newNode(rootNodeX-SUNS_SOLAR_SYSTEM_RADIUS*zLevel, rootNodeY))
	results = append(results, newNode(rootNodeX+SUNS_SOLAR_SYSTEM_RADIUS*zLevel, rootNodeY))
	results = append(results, newNode(rootNodeX-horizontalOffset*zLevel, rootNodeY+verticalOffset*zLevel))
	results = append(results, newNode(rootNodeX-horizontalOffset*zLevel, rootNodeY-verticalOffset*zLevel))
	results = append(results, newNode(rootNodeX+horizontalOffset*zLevel, rootNodeY+verticalOffset*zLevel))
	results = append(results, newNode(rootNodeX+horizontalOffset*zLevel, rootNodeY-verticalOffset*zLevel))

	for i := 1; i < zLevel; i++ {
		results = append(results, newNode(rootNodeX-SUNS_SOLAR_SYSTEM_RADIUS*zLevel+horizontalOffset*i, rootNodeY+verticalOffset*i))
		results = append(results, newNode(rootNodeX-SUNS_SOLAR_SYSTEM_RADIUS*zLevel+horizontalOffset*i, rootNodeY-verticalOffset*i))
		results = append(results, newNode(rootNodeX+SUNS_SOLAR_SYSTEM_RADIUS*zLevel-horizontalOffset*i, rootNodeY+verticalOffset*i))
		results = append(results, newNode(rootNodeX+SUNS_SOLAR_SYSTEM_RADIUS*zLevel-horizontalOffset*i, rootNodeY-verticalOffset*i))
		results = append(results, newNode(rootNodeX-horizontalOffset*zLevel+SUNS_SOLAR_SYSTEM_RADIUS*i, rootNodeY+verticalOffset*zLevel))
		results = append(results, newNode(rootNodeX-horizontalOffset*zLevel+SUNS_SOLAR_SYSTEM_RADIUS*i, rootNodeY-verticalOffset*zLevel))
	}
	return results
}

type Node struct {
	Data     string
	Position *vec2d.Vector
}

func newNode(x, y float64) *Node {
	n := new(Node)
	n.Position = vec2d.New(x, y)
	return n
}

func (n *Node) Key() string {
	return fmt.Sprintf("Node:%d_%d", n.Position.X, n.Position.Y)
}

// TODO: This node thing should be a type
func (s *Sun) findHomeNode(rootNode *Node) *Node {
	var zLevel uint32 = 1

	rootNodeEntity, _ := Get(rootNode.Key())

	if rootNodeEntity.(*Node).Data == "" {
		return rootNode
	}

	for {
		nodes := fetchNodesLayer(rootNodeEntity.(*Node), zLevel)
		nodesEntities := FindList(nodes)
		for _, nodeEntity := range nodesEntities {
			if nodeEntity.Data == "" {
				return nodeEntity
			}
		}
		zLevel++
	}
}

// Uses all player's twitter friends and tries to place the sun as
// close as possible to all of them. This of course could cause tons of
// overlapping. To solve this, we simply throw the sun somewhere far away
// from the desired point and start to move it to THE POINT, but carefully
// watching for collisions.
func GenerateSun(username string, friends, others []Sun) *Sun {
	newSun := Sun{
		Username: username,
		speed:    5,
		target:   vec2d.New(0, 0),
		Position: getRandomStartPosition(SUNS_RANDOM_SPAWN_ZONE_RADIUS),
	}
	newSun.generateName(username)

	node := findHomeNode(getStartPosNode(friends))
	node.Data = newSun.Key()
	Save(node)

	newSun.Position.X = node.Position.X
	newSun.Position.Y = node.Position.Y
	return &newSun
}
