/*
 * Generated by asn1c-0.9.29 (http://lionet.info/asn1c)
 * From ASN.1 module "DSRC"
 * 	found in "j2735.asn"
 * 	`asn1c -fcompound-names -pdu=auto`
 */

#ifndef	_TestMessage04_H_
#define	_TestMessage04_H_


#include <asn_application.h>

/* Including external dependencies */
#include <constr_SEQUENCE.h>

#ifdef __cplusplus
extern "C" {
#endif

/* Forward declarations */
struct Header;
struct RegionalExtension;

/* TestMessage04 */
typedef struct TestMessage04 {
	struct Header	*header	/* OPTIONAL */;
	struct RegionalExtension	*regional	/* OPTIONAL */;
	/*
	 * This type is extensible,
	 * possible extensions are below.
	 */
	
	/* Context for parsing across buffer boundaries */
	asn_struct_ctx_t _asn_ctx;
} TestMessage04_t;

/* Implementation */
extern asn_TYPE_descriptor_t asn_DEF_TestMessage04;

#ifdef __cplusplus
}
#endif

/* Referred external types */
#include "Header.h"
#include "RegionalExtension.h"

#endif	/* _TestMessage04_H_ */
#include <asn_internal.h>
