// package dotcommonitor

// import (
// 	"testing"
// )

// Common structures

// Device structures
// func TestExpandNotificationsNotificationGroupList(t *testing.T) {
// 	cases := []struct {
//         ownerId  *string
//         pairs    []*ec2.UserIdGroupPair
//         expected []*GroupIdentifier
//     }{
//         // simple, no user id included
//         {
//             ownerId: aws.String("user1234"),
//             pairs: []*ec2.UserIdGroupPair{
//                 &ec2.UserIdGroupPair{
//                     GroupId: aws.String("sg-12345"),
//                 },
//             },
//             expected: []*GroupIdentifier{
//                 &GroupIdentifier{
//                     GroupId: aws.String("sg-12345"),
//                 },
//             },
//         },
// 	}
//}