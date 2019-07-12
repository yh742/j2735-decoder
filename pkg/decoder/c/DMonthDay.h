/*
 * Generated by asn1c-0.9.29 (http://lionet.info/asn1c)
 * From ASN.1 module "DSRC"
 * 	found in "j2735.asn"
 * 	`asn1c -fcompound-names -pdu=auto`
 */

#ifndef	_DMonthDay_H_
#define	_DMonthDay_H_


#include <asn_application.h>

/* Including external dependencies */
#include "DMonth.h"
#include "DDay.h"
#include <constr_SEQUENCE.h>

#ifdef __cplusplus
extern "C" {
#endif

/* DMonthDay */
typedef struct DMonthDay {
	DMonth_t	 month;
	DDay_t	 day;
	
	/* Context for parsing across buffer boundaries */
	asn_struct_ctx_t _asn_ctx;
} DMonthDay_t;

/* Implementation */
extern asn_TYPE_descriptor_t asn_DEF_DMonthDay;

#ifdef __cplusplus
}
#endif

#endif	/* _DMonthDay_H_ */
#include <asn_internal.h>
