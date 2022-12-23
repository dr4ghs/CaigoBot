package join

const (
	// Only allows the owner to accept join requests
	JoinRequestAcceptPolicyOwner JoinRequestAcceptPolicy = iota

	// Owner and admins can accept join requests
	JoinRequestAcceptPolicyAdmins

	// Owner, admins and moderators can accept join requests
	JoinRequestAcceptPolicyModerators

	// All owner's staff members can accept join requests
	JoinRequestAcceptPolicyStaff
)

const (
	// Only owner's staff members can join the stream room
	JoinRequestJoinPolicyStaff JoinRequestJoinPolicy = iota

	// Only owner's staff members and subscribers can join the stream room
	JoinRequestJoinPolicySubscribers

	// Only owner's staff members, subscribers and followers can join the stream room
	JoinRequestJoinPolicyFollowers

	// Everyone can join the stream room
	JoinRequestJoinPolicyEveryone
)
