package tables

import (
	"fmt"
	"github.com/80asis/cyclops/entity"
)

// Define an interface for LocalTable
type LocalTableInterface interface {
	FetchConnectedAZs() []string
	FetchAllEntities() []entity.Entity
	VerifyChecksum(entityID, checksum string) bool
	UpdateChecksum(entityID, checksum string)
	MarkOutOfSync(entityID string)
	UpdateSyncStatus(entityID string, status string)
}

// Define an interface for RemoteTable
type RemoteTableInterface interface {
	VerifyChecksumInAZ(entityID, checksum, az string) bool
	UpdateSyncStatus(entityID string, status string)
}

// LocalTable represents the local database table
type LocalTable struct{}

// FetchConnectedAZs fetches all connected AZs and returns their FQDN
func (table *LocalTable) FetchConnectedAZs() []string {
	// Implement logic to fetch connected AZs
	return []string{"AZ1", "AZ2"} // Dummy data
}

// Fetches all the entities that are present in local AZ for sync
func (table *LocalTable) FetchAllEntities() []entity.Entity {
	// Implement logic to fetch connected AZs
	return []entity.Entity{
		{EntityID: "entity1", EntityKind: "protection_rule", OpType: "update"},
		{EntityID: "entity2", EntityKind: "recovery_plan", OpType: "delete"},
		{EntityID: "entity3", EntityKind: "category", OpType: "update"},
	} // Dummy data
}

// VerifyChecksum verifies the checksum of an entity in the local table
func (table *LocalTable) VerifyChecksum(entityID, checksum string) bool {
	// Implement logic to verify checksum
	// Example: Retrieve checksum from local table for the entityID and compare with the provided checksum
	return false // Dummy data
}

// UpdateChecksum updates the checksum of an entity in the local table
func (table *LocalTable) UpdateChecksum(entityID, checksum string) {
	// Implement logic to update checksum
	fmt.Printf("Updated checksum for entityID: %s\n", entityID)
}

// MarkOutOfSync marks an entity as out of sync in the local table
func (table *LocalTable) MarkOutOfSync(entityID string) {
	// Implement logic to mark entity as out of sync
	fmt.Printf("Marked entity as out of sync: %s\n", entityID)
}

func (table *LocalTable) UpdateSyncStatus(entityID string, status string) {
	fmt.Printf("Sync status %s updated\n", status)
}

// RemoteTable represents the remote database table
type RemoteTable struct{}

// VerifyChecksumInAZ verifies the checksum of an entity in the remote table for a specific AZ
func (table *RemoteTable) VerifyChecksumInAZ(entityID, checksum, az string) bool {
	// Implement logic to verify checksum in remote table for a specific AZ
	// Example: Retrieve checksum from remote table for the entityID in the specified AZ and compare with the provided checksum
	return false // Dummy data
}

func (table *RemoteTable) UpdateSyncStatus(entityID string, status string) {
	fmt.Printf("Sync status %s updated\n", status)
}
