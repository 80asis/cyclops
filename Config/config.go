package Config

import (
	"github.com/80asis/cyclops/cyclops"
)

// # The list always needs to be ordered such that the dependee entity/plugin is
// # listed before the dependent entity/plugin. This list should be restricted to
// # only global entities. Any entities that the global entities depend on (eg:
// # category for ProtectionRule) should be taken care by the global entity
// # itself.
var RegisterPlugins = []cyclops.EntitySyncEntityType_Type{
	cyclops.EntitySyncEntityType_kCategory,
	cyclops.EntitySyncEntityType_kProtectionRule,
}

// # Mapping from EntitySyncEntityType enum to Entity type string. This will be
// # shown under entity type on UI.
var Entity_type_str_map = map[cyclops.EntitySyncEntityType_Type]string{
	cyclops.EntitySyncEntityType_kCategory:       "category",
	cyclops.EntitySyncEntityType_kProtectionRule: "protection_rule",
}

// # Mapping of entity types to their corresponding class names.
var Plugin_obj_mapping = map[cyclops.EntitySyncEntityType_Type]interface{}{
	cyclops.EntitySyncEntityType_kCategory:       "category",        // this should have category interface here
	cyclops.EntitySyncEntityType_kProtectionRule: "protection_rule", // thsi should have protetctionRule interface here
}
