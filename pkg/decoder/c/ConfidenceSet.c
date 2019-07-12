/*
 * Generated by asn1c-0.9.29 (http://lionet.info/asn1c)
 * From ASN.1 module "DSRC"
 * 	found in "j2735.asn"
 * 	`asn1c -fcompound-names -pdu=auto`
 */

#include "ConfidenceSet.h"

asn_TYPE_member_t asn_MBR_ConfidenceSet_1[] = {
	{ ATF_POINTER, 7, offsetof(struct ConfidenceSet, accelConfidence),
		(ASN_TAG_CLASS_CONTEXT | (0 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_AccelSteerYawRateConfidence,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"accelConfidence"
		},
	{ ATF_POINTER, 6, offsetof(struct ConfidenceSet, speedConfidence),
		(ASN_TAG_CLASS_CONTEXT | (1 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_SpeedandHeadingandThrottleConfidence,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"speedConfidence"
		},
	{ ATF_POINTER, 5, offsetof(struct ConfidenceSet, timeConfidence),
		(ASN_TAG_CLASS_CONTEXT | (2 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_TimeConfidence,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"timeConfidence"
		},
	{ ATF_POINTER, 4, offsetof(struct ConfidenceSet, posConfidence),
		(ASN_TAG_CLASS_CONTEXT | (3 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_PositionConfidenceSet,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"posConfidence"
		},
	{ ATF_POINTER, 3, offsetof(struct ConfidenceSet, steerConfidence),
		(ASN_TAG_CLASS_CONTEXT | (4 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_SteeringWheelAngleConfidence,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"steerConfidence"
		},
	{ ATF_POINTER, 2, offsetof(struct ConfidenceSet, headingConfidence),
		(ASN_TAG_CLASS_CONTEXT | (5 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_HeadingConfidence,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"headingConfidence"
		},
	{ ATF_POINTER, 1, offsetof(struct ConfidenceSet, throttleConfidence),
		(ASN_TAG_CLASS_CONTEXT | (6 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_ThrottleConfidence,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"throttleConfidence"
		},
};
static const int asn_MAP_ConfidenceSet_oms_1[] = { 0, 1, 2, 3, 4, 5, 6 };
static const ber_tlv_tag_t asn_DEF_ConfidenceSet_tags_1[] = {
	(ASN_TAG_CLASS_UNIVERSAL | (16 << 2))
};
static const asn_TYPE_tag2member_t asn_MAP_ConfidenceSet_tag2el_1[] = {
    { (ASN_TAG_CLASS_CONTEXT | (0 << 2)), 0, 0, 0 }, /* accelConfidence */
    { (ASN_TAG_CLASS_CONTEXT | (1 << 2)), 1, 0, 0 }, /* speedConfidence */
    { (ASN_TAG_CLASS_CONTEXT | (2 << 2)), 2, 0, 0 }, /* timeConfidence */
    { (ASN_TAG_CLASS_CONTEXT | (3 << 2)), 3, 0, 0 }, /* posConfidence */
    { (ASN_TAG_CLASS_CONTEXT | (4 << 2)), 4, 0, 0 }, /* steerConfidence */
    { (ASN_TAG_CLASS_CONTEXT | (5 << 2)), 5, 0, 0 }, /* headingConfidence */
    { (ASN_TAG_CLASS_CONTEXT | (6 << 2)), 6, 0, 0 } /* throttleConfidence */
};
asn_SEQUENCE_specifics_t asn_SPC_ConfidenceSet_specs_1 = {
	sizeof(struct ConfidenceSet),
	offsetof(struct ConfidenceSet, _asn_ctx),
	asn_MAP_ConfidenceSet_tag2el_1,
	7,	/* Count of tags in the map */
	asn_MAP_ConfidenceSet_oms_1,	/* Optional members */
	7, 0,	/* Root/Additions */
	7,	/* First extension addition */
};
asn_TYPE_descriptor_t asn_DEF_ConfidenceSet = {
	"ConfidenceSet",
	"ConfidenceSet",
	&asn_OP_SEQUENCE,
	asn_DEF_ConfidenceSet_tags_1,
	sizeof(asn_DEF_ConfidenceSet_tags_1)
		/sizeof(asn_DEF_ConfidenceSet_tags_1[0]), /* 1 */
	asn_DEF_ConfidenceSet_tags_1,	/* Same as above */
	sizeof(asn_DEF_ConfidenceSet_tags_1)
		/sizeof(asn_DEF_ConfidenceSet_tags_1[0]), /* 1 */
	{ 0, 0, SEQUENCE_constraint },
	asn_MBR_ConfidenceSet_1,
	7,	/* Elements count */
	&asn_SPC_ConfidenceSet_specs_1	/* Additional specs */
};

