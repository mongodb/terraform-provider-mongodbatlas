package acc

import (
	"fmt"
	"sync/atomic"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

const (
	prefixName     = "test-acc-tf"
	prefixProject  = prefixName + "-p"
	prefixCluster  = prefixName + "-c"
	prefixStream   = prefixName + "-s"
	prefixIAMRole  = "mongodb-atlas-" + prefixName
	prefixIAMUser  = "arn:aws:iam::358363220050:user/mongodb-aws-iam-auth-test-user"
	prefixS3Bucket = "mongodb-atlas-tf"
)

func RandomName() string {
	return acctest.RandomWithPrefix(prefixName)
}

func RandomProjectName() string {
	return acctest.RandomWithPrefix(prefixProject)
}

func RandomClusterName() string {
	return acctest.RandomWithPrefix(prefixCluster)
}

func RandomStreamInstanceName() string {
	return acctest.RandomWithPrefix(prefixStream)
}

func RandomIAMRole() string {
	return acctest.RandomWithPrefix(prefixIAMRole)
}

func RandomIAMUser() string {
	return acctest.RandomWithPrefix(prefixIAMUser)
}

// ipCounter is a process-global counter used by RandomIP to allocate unique IPs.
// Initialized to a random offset so two test runs do not replay identical sequences.
var ipCounter atomic.Uint32

const (
	// ipCounterMask keeps the counter inside a 24-bit address space (3 octets).
	ipCounterMask uint32 = 0xFFFFFF
	// ipOctetMask isolates a single octet (8 bits) from the counter.
	ipOctetMask uint32 = 0xFF
	// ipOctetShiftHigh / ipOctetShiftMid extract the high and middle octets
	// from the 24-bit counter; the low octet is taken without a shift.
	ipOctetShiftHigh uint32 = 16
	ipOctetShiftMid  uint32 = 8
)

func init() {
	ipCounter.Store(uint32(acctest.RandIntRange(0, int(ipCounterMask)))) //nolint:gosec // value capped at 24 bits, fits in uint32
}

// RandomIP returns a unique IP within the 179.0.0.0/8 range. Uniqueness is
// guaranteed within a single test process for up to ~16M calls, which avoids
// the collisions a previous shared-octets implementation produced when parallel
// tests rolled the same last octet.
func RandomIP() string {
	n := ipCounter.Add(1) & ipCounterMask
	return fmt.Sprintf("179.%d.%d.%d",
		byte((n>>ipOctetShiftHigh)&ipOctetMask),
		byte((n>>ipOctetShiftMid)&ipOctetMask),
		byte(n&ipOctetMask),
	)
}

func RandomEmail() string {
	return fmt.Sprintf("%s-%s@mongodb.com", prefixName, acctest.RandString(10))
}

func RandomLDAPName() string {
	return fmt.Sprintf("CN=%s-%s@example.com,OU=users,DC=example,DC=com", prefixName, acctest.RandString(10))
}

func RandomBucketName() string {
	return fmt.Sprintf("%s-%s", prefixS3Bucket, acctest.RandString(10))
}
