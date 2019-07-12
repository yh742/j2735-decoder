/*
 * Generated by asn1c-0.9.29 (http://lionet.info/asn1c)
 * From ASN.1 module "DSRC"
 * 	found in "j2735.asn"
 * 	`asn1c -fcompound-names -pdu=auto`
 */

#include "SpecialVehicleExtensions.h"

asn_TYPE_member_t asn_MBR_SpecialVehicleExtensions_1[] = {
	{ ATF_POINTER, 3, offsetof(struct SpecialVehicleExtensions, vehicleAlerts),
		(ASN_TAG_CLASS_CONTEXT | (0 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_EmergencyDetails,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"vehicleAlerts"
		},
	{ ATF_POINTER, 2, offsetof(struct SpecialVehicleExtensions, description),
		(ASN_TAG_CLASS_CONTEXT | (1 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_EventDescription,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"description"
		},
	{ ATF_POINTER, 1, offsetof(struct SpecialVehicleExtensions, trailers),
		(ASN_TAG_CLASS_CONTEXT | (2 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_TrailerData,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"trailers"
		},
};
static const int asn_MAP_SpecialVehicleExtensions_oms_1[] = { 0, 1, 2 };
static const ber_tlv_tag_t asn_DEF_SpecialVehicleExtensions_tags_1[] = {
	(ASN_TAG_CLASS_UNIVERSAL | (16 << 2))
};
static const asn_TYPE_tag2member_t asn_MAP_SpecialVehicleExtensions_tag2el_1[] = {
    { (ASN_TAG_CLASS_CONTEXT | (0 << 2)), 0, 0, 0 }, /* vehicleAlerts */
    { (ASN_TAG_CLASS_CONTEXT | (1 << 2)), 1, 0, 0 }, /* description */
    { (ASN_TAG_CLASS_CONTEXT | (2 << 2)), 2, 0, 0 } /* trailers */
};
asn_SEQUENCE_specifics_t asn_SPC_SpecialVehicleExtensions_specs_1 = {
	sizeof(struct SpecialVehicleExtensions),
	offsetof(struct SpecialVehicleExtensions, _asn_ctx),
	asn_MAP_SpecialVehicleExtensions_tag2el_1,
	3,	/* Count of tags in the map */
	asn_MAP_SpecialVehicleExtensions_oms_1,	/* Optional members */
	3, 0,	/* Root/Additions */
	3,	/* First extension addition */
};
asn_TYPE_descriptor_t asn_DEF_SpecialVehicleExtensions = {
	"SpecialVehicleExtensions",
	"SpecialVehicleExtensions",
	&asn_OP_SEQUENCE,
	asn_DEF_SpecialVehicleExtensions_tags_1,
	sizeof(asn_DEF_SpecialVehicleExtensions_tags_1)
		/sizeof(asn_DEF_SpecialVehicleExtensions_tags_1[0]), /* 1 */
	asn_DEF_SpecialVehicleExtensions_tags_1,	/* Same as above */
	sizeof(asn_DEF_SpecialVehicleExtensions_tags_1)
		/sizeof(asn_DEF_SpecialVehicleExtensions_tags_1[0]), /* 1 */
	{ 0, 0, SEQUENCE_constraint },
	asn_MBR_SpecialVehicleExtensions_1,
	3,	/* Elements count */
	&asn_SPC_SpecialVehicleExtensions_specs_1	/* Additional specs */
};

