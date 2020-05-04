package domain

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// --------------------  Time tracking related types ------------------

//TimeTrackedEntity is an interface that is obeyed
//from all objects that have a time dimension, meaning
//that come into existence in some specific point in time (pit)
//and may stop being in a later pit. A TimeTrackedEntity
//cannot be "stopped" in a pit and then later come again
//"alive"
type TimeTrackedEntity interface {

	//IsExistentAt returns true if the
	//current object is still valid,
	//meaning that it doesn't have
	//an ending time (and has started being in a pit)
	IsExistentAt(pit time.Time) bool

	//ExistentFrom returns the time this object
	//started to exist. A TimedTrackedEntity SHOULD ALWAYS
	//have a started pit
	ExistentFrom() time.Time

	//ValidUntil returns the pit that this
	//object came out of existence.
	//If this object is still "active" then
	//calling this method SHOULD return NilTime
	ValidUntil() time.Time

	//ActiveDuration returns the duration
	//that this entity is "existent". If
	//the entity is still in an active state
	//then calling this method returns the
	//difference between activation and now
	ActiveDuration() time.Duration
}

//------------------------------------------------------------------

//TimeTrackedEntityCollection is a structure used
//for storing time track entities and provide
//utility functions for searching and filtering
//them.
//It is based on an augmented interval tree
type TimeTrackedEntityCollection struct {
	root      *intervalNode
	noOfNodes int
}

//String implementation traverse the collection and
//return the result
func (ts TimeTrackedEntityCollection) String() string {

	var str strings.Builder

	ts.traverseNodes(ts.root, func(n *intervalNode, level int) {
		str.WriteString("(" + strconv.Itoa(level) + ")" + n.String())
	}, 0)

	return str.String()
}

//AddEntity adds a new entity to the tracked
//collections. It doesn't test if the entity
//already exists in the collection
func (ts *TimeTrackedEntityCollection) AddEntity(e TimeTrackedEntity) {

	newNodeToInsert := &intervalNode{
		entity: e,
		max:    e.ValidUntil(),
		left:   nil,
		right:  nil,
	}

	ts.root = ts.insertNode(ts.root, newNodeToInsert)
	ts.noOfNodes++
}

func (ts *TimeTrackedEntityCollection) intersectNode(tmp *intervalNode, searchFor TimeTrackedEntity, foundSoFar []TimeTrackedEntity) {

	if tmp == nil {
		return
	}

	if !searchFor.ValidUntil().IsZero() {
		if !(compareEndTime(tmp.entity.ExistentFrom(), searchFor.ValidUntil()) < 0 ||
			compareEndTime(tmp.entity.ValidUntil(), searchFor.ExistentFrom()) > 0) {

		}
	}

}

//InsertEntity adds an entity to the collections
func (ts *TimeTrackedEntityCollection) insertNode(tmp *intervalNode, newNode *intervalNode) *intervalNode {

	// Check if we are in
	if tmp == nil {
		return newNode
	}

	//Check to see if the newly added node
	//has and ending that this further the current max
	//for this node
	if compareEndTime(tmp.max, newNode.max) < 0 {
		tmp.max = newNode.max
	}

	//proceed with insertion
	if tmp.compareTo(newNode) <= 0 {
		if tmp.right == nil {
			tmp.right = newNode
		} else {
			ts.insertNode(tmp.right, newNode)
		}
	} else {
		if tmp.left == nil {
			tmp.left = newNode
		} else {
			ts.insertNode(tmp.left, newNode)
		}
	}
	return tmp
}

// visitorFunc is a function
// that is used when visiting a node
// of a TimeTrackedEntityCollection
type visitorFunc func(n *intervalNode, level int)

//traverseNodes , performa a pre order traversal and calls visitor
//in every node visited
func (ts *TimeTrackedEntityCollection) traverseNodes(n *intervalNode, visitor visitorFunc, currentLevel int) {

	if n == nil {
		return
	}

	// visit left subtree
	if n.left != nil {
		ts.traverseNodes(n.left, visitor, currentLevel+1)
	}

	// visit n
	visitor(n, currentLevel)

	//visit right sub tree
	if n.right != nil {
		ts.traverseNodes(n.right, visitor, currentLevel+1)
	}

}

//-----------------------------------------------------------
//                   Utility functions
//-----------------------------------------------------------

//NilTime is a utility function that returns
//the equivalent for nil for time.Time
//Someone can test if null with
func NilTime() time.Time {
	return time.Time{}
}

// Compares two ending times , taking into account
// the concept of "not ended yet".
// Returns -1 if a ends before b
// Returns 1 if a ends after b
// Return 0 if a and b are equal (includes the case for Zero Time)
func compareEndTime(a time.Time, b time.Time) int {

	if a.IsZero() {
		if b.IsZero() {
			return 0
		}
		// b ends before a (which has no ending time)
		return 1
	}

	if b.IsZero() {
		//a has a defined ending, b has not
		return -1
	}

	if a.Before(b) {
		return -1
	}

	if a == b {
		return 0
	}

	return 1

}

// ------------------------------------------------

//intervalNode is a concrete augmented node
//of the interval tree
type intervalNode struct {
	// the entity that is kept in the node
	entity TimeTrackedEntity
	// the maximum ending time of the
	// tree below this node
	max time.Time
	// left subtree
	left *intervalNode
	// right subtree
	right *intervalNode
}

//compareTo , compares a node with another. The comparison
//is based on the starting point of the contained time tracked
//entity and subsequently to the duration of the entity.
//Returns -1 if n starts before the compare to node and 1 otherwise.
//If they are equal it retuns 0
func (n intervalNode) compareTo(anotherNode *intervalNode) int {

	if n.entity.ExistentFrom().Before(anotherNode.entity.ExistentFrom()) {
		return -1
	} else if n.entity.ExistentFrom() == anotherNode.entity.ExistentFrom() {
		return compareEndTime(n.entity.ValidUntil(), anotherNode.entity.ValidUntil())
	}

	return 1
}

//String implementation of a node
func (n intervalNode) String() string {
	return fmt.Sprintf("[E:%v M:%v]", n.entity, n.max)
}

//------------------------------------------------

//------------ Dynamic Attributes related types ----

//AttributeBearer is an interface that when fullfilled
//from an object, dynamic attributes can be assigned to this
//object
type AttributeBearer interface {

	//GetAttributeNames return the name
	//of all the attributes the current
	//object has
	GetAttributeNames() []string

	//HasAttribute checks if an attribute is present
	//in the current entity.
	HasAttribute(attrName string) bool

	//GetAttribute returns the value of the attribute
	//or an error if this atrribute did not exists
	GetAttribute(attrName string) (interface{}, error)

	//SetAttribute set the value for a given attribute.
	//If the attribute already exists then it is overriden
	//and the previous value is returned Otherwise is added
	//and nil is returned
	SetAttribute(attrName string, value interface{}) interface{}
}
