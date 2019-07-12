/*
 * Generated by asn1c-0.9.29 (http://lionet.info/asn1c)
 * From ASN.1 module "AddGrpC"
 * 	found in "j2735.asn"
 * 	`asn1c -fcompound-names -pdu=auto`
 */

#include "IntersectionState-addGrpC.h"

static asn_TYPE_member_t asn_MBR_IntersectionState_addGrpC_1[] = {
	{ ATF_POINTER, 1, offsetof(struct IntersectionState_addGrpC, activePrioritizations),
		(ASN_TAG_CLASS_CONTEXT | (0 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_PrioritizationResponseList,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"activePrioritizations"
		},
};
static const int asn_MAP_IntersectionState_addGrpC_oms_1[] = { 0 };
static const ber_tlv_tag_t asn_DEF_IntersectionState_addGrpC_tags_1[] = {
	(ASN_TAG_CLASS_UNIVERSAL | (16 << 2))
};
static const asn_TYPE_tag2member_t asn_MAP_IntersectionState_addGrpC_tag2el_1[] = {
    { (ASN_TAG_CLASS_CONTEXT | (0 << 2)), 0, 0, 0 } /* activePrioritizations */
};
static asn_SEQUENCE_specifics_t asn_SPC_IntersectionState_addGrpC_specs_1 = {
	sizeof(struct IntersectionState_addGrpC),
	offsetof(struct IntersectionState_addGrpC, _asn_ctx),
	asn_MAP_IntersectionState_addGrpC_tag2el_1,
	1,	/* Count of tags in the map */
	asn_MAP_IntersectionState_addGrpC_oms_1,	/* Optional members */
	1, 0,	/* Root/Additions */
	1,	/* First extension addition */
};
asn_TYPE_descriptor_t asn_DEF_IntersectionState_addGrpC = {
	"IntersectionState-addGrpC",
	"IntersectionState-addGrpC",
	&asn_OP_SEQUENCE,
	asn_DEF_IntersectionState_addGrpC_tags_1,
	sizeof(asn_DEF_IntersectionState_addGrpC_tags_1)
		/sizeof(asn_DEF_IntersectionState_addGrpC_tags_1[0]), /* 1 */
	asn_DEF_IntersectionState_addGrpC_tags_1,	/* Same as above */
	sizeof(asn_DEF_IntersectionState_addGrpC_tags_1)
		/sizeof(asn_DEF_IntersectionState_addGrpC_tags_1[0]), /* 1 */
	{ 0, 0, SEQUENCE_constraint },
	asn_MBR_IntersectionState_addGrpC_1,
	1,	/* Elements count */
	&asn_SPC_IntersectionState_addGrpC_specs_1	/* Additional specs */
};

