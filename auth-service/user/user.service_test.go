package user_test

import (
	"context"
	"encoding/base64"
	"os"
	"testing"
	"time"

	"log/slog"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/sing3demons/auth-service/redis"
	"github.com/sing3demons/auth-service/store"
	"github.com/sing3demons/auth-service/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const PUBLIC_ACCESS_KEY = "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQ0lUQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FnNEFNSUlDQ1FLQ0FnQjFOUE1pWmxSSjVIMk00cmV1aGFiQgpXaDliZU1VQi8vWXBPNmVoSmtiTzd2ZTdyTWh0ZC9xRHFXSWd6cGdtNm0yL0lMTUJCck5CbzZWUnhqVUtHSTBMCk1PVk81a09pejZWc3BBRDd6RlFNOXRSbXZVOWdrRE85U0dsdHhxaHZxS3FTdHYzNWZ3blJzNlhERGJTLzVLMFgKdWZBeWY3NlBUQkZ6NnBRZnROU1lhRURGeVkrM2s2eGlDTlB5Vkw1Y05LeWlPQzVBNllMb3FWZHhHV2RGYkFVVAo2Z3pseldhNXhHUm9ZZXZnaFA4N01yMDNSNTdlRVk5bnV3cXpMc2lpWWxUd0JIOGMvNXQxdWFyZlNEMnBFajZyCmlzUytMSGkzc0h3RFRnRFE4UEt5ZEt3bytmNzhNN2s0VDR3bld6Nlp0ZC84UFFYZ1d6dm9pWk5BZzhJbnBpRTUKUmlYT3dMcHlraS9YYUUvNlFzTkM1TjhYZVVIUDRta0UvUjFuSHFRaTBVNXpKbzVFUnhQRzNVeHkyYVI3US9ZZgptZHNlUUt5WEtETHIxSUk0Z25DcjNlQmxGLzJrT2Z2NWszZXVxMS92c1l2S1k2NThEb0U4TFBBK2t5QkhFQTVRCkp4RkkxZlVWSTVwWmtzdnJPeVFGUURqSjFXN2Y3UG5EOTB5WnpWSnIrYWxhMGV5eWdDcjNoSGpjZGNvZEpXQkEKL3dCUzhQbGxEclBmbWdUWkRQTHZLcWNTcGh3WGRXZG92aEpFZk44L3dBQUxJNmlBQWgzVnJORmJ3NkZZdHVvTgppcEMwYStNczVlVnMzL1duU3ZtRy8zTFRSZkh0VXRYRDZZMkNxV3pjRi9GeHZCV2lwVlB6YlAxblllM0pObU53CkdHQ2JpR3ZJUkNiQWhXSjFQRVRKWVFJREFRQUIKLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0t"
const PUBLIC_REFRESH_KEY = "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQ0lqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FnOEFNSUlDQ2dLQ0FnRUFpR2FNRDZBL3JpYlZQS1U3N0pSbApuejE3V0tsNmVzVVhpcDFjL3RVQVNuK3FYMlhrL1dWZmtBaWgrZDJKV2hkR1BCQkZoZ2JNdXBMZTNHTDJjYXZsCnBaT2wwanMvTW9zNVlVU3BZcVpPVnJydGlQb08vU1U1YzFKZk9oekQwYncvUmhrSXc4ay8vMkQzTXlORkh4aHgKOVhhY2Z0aTIzaEFzWk5rUjYwcEtjaWpvcFFFSEcxYjBUaGlaK0lKdy92RktYbEcyNFhLKzQrUWRyOUJrbndOZQpxbkJ0VWZHOXpYaUdRaXN2ZEJLdkNwUi9rdi9KSFNWa3N6NkM0Qk9QS25UVXc4aWNHbitXNHZRSm1qbmxyOGhyClVVcm16bGM0UW02T2UvR2RGcmJNVkdqcE1rQ0grSE1CMVNkMGsxaithV2xRMTFuWGJ1TDZOUThGNmRjdVVqWWEKbE8zTm41ZG83L3NvWmFNV3dYU2xSYjRDWFJ0QlN3ZmU4Vmt0dFU0OEpSSFJ6ZE5GZ0pwZTVCWUdHZGJOb0tHRQpmNW1OWkY3VU81YWZnUENMQWk0bnY5aGNrWVI2Rmd3eVpTcTk2c0FLT2pPaHp6SzdDNVNDQXlPTXNkSDR1QlNpCjM0djZtVFNYdDhpWGoyWHNsaVBqU3BPaTBBT29ldEVudERjZU9xWXlXMzZzRDRGNHJ1YkpTWXU2VUZBT2JaWlkKWkppckp0Vnc0Z3diNDFQZzdWbUFUSkh6alI0cmQxaDVjcWQ2UzlaKzI5a2dUTFBMQkYwZlJNTDVTOElKV3ptQgpuc0g5VTdzNlMxeEVueWVPczAySTlXa2VwOVFkZmlxQnVLanV3cUJuVlVUR284ZkdUSmtlbzBKZUoxbjFxMTl1CklYWC82TnVsZTNZcFhHTGxwcFVUbWZNQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ=="
const PRIVATE_ACCESS_KEY = "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlKSmdJQkFBS0NBZ0IxTlBNaVpsUko1SDJNNHJldWhhYkJXaDliZU1VQi8vWXBPNmVoSmtiTzd2ZTdyTWh0CmQvcURxV0lnenBnbTZtMi9JTE1CQnJOQm82VlJ4alVLR0kwTE1PVk81a09pejZWc3BBRDd6RlFNOXRSbXZVOWcKa0RPOVNHbHR4cWh2cUtxU3R2MzVmd25SczZYRERiUy81SzBYdWZBeWY3NlBUQkZ6NnBRZnROU1lhRURGeVkrMwprNnhpQ05QeVZMNWNOS3lpT0M1QTZZTG9xVmR4R1dkRmJBVVQ2Z3pseldhNXhHUm9ZZXZnaFA4N01yMDNSNTdlCkVZOW51d3F6THNpaVlsVHdCSDhjLzV0MXVhcmZTRDJwRWo2cmlzUytMSGkzc0h3RFRnRFE4UEt5ZEt3bytmNzgKTTdrNFQ0d25XejZadGQvOFBRWGdXenZvaVpOQWc4SW5waUU1UmlYT3dMcHlraS9YYUUvNlFzTkM1TjhYZVVIUAo0bWtFL1IxbkhxUWkwVTV6Sm81RVJ4UEczVXh5MmFSN1EvWWZtZHNlUUt5WEtETHIxSUk0Z25DcjNlQmxGLzJrCk9mdjVrM2V1cTEvdnNZdktZNjU4RG9FOExQQStreUJIRUE1UUp4RkkxZlVWSTVwWmtzdnJPeVFGUURqSjFXN2YKN1BuRDkweVp6VkpyK2FsYTBleXlnQ3IzaEhqY2Rjb2RKV0JBL3dCUzhQbGxEclBmbWdUWkRQTHZLcWNTcGh3WApkV2RvdmhKRWZOOC93QUFMSTZpQUFoM1ZyTkZidzZGWXR1b05pcEMwYStNczVlVnMzL1duU3ZtRy8zTFRSZkh0ClV0WEQ2WTJDcVd6Y0YvRnh2QldpcFZQemJQMW5ZZTNKTm1Od0dHQ2JpR3ZJUkNiQWhXSjFQRVRKWVFJREFRQUIKQW9JQ0FDVEw1ZFVUNlR4MWxwRUhrSUVqQnBKSFYvYmd1SUVET2VZQ0M0T0ZQOCt4cUdic1BOUlpTWFhTVkxOVwpDT0NXMHJPaGNYRk9DRE1BVEdPYTVZWHc1VDd4TDl5UVlBV2FTU1lOYXgyaUxYVVFmT2h3WUo1QlIyMFNjYjc4CkVsOVR4WkZnRCtZbll3N0o1cTJRL1FFTnF1WDdBeFRubEF6cTVjUE5qb2xSdlRqSDZpWHVQTWQyZmpVYzdtVnoKQTN4eE1RMlFzN3kvVXREMmNUUlp0RmxRSzF2d0svSnRoT1duYVpwM3U1VDJUQ2JxckFyUjJtZC9mZFVrSEp6YQpnRFN5eUZXK0k2WVgvVmQ5WGM3Um9FSVhMME90d2t5dUs0d0h4VGRJWXVzMTZndnFveTRPSm1aSldzbjU5MWRYCnh4Ujh6QUFsckZXZWJ2VjNVNXA2Z3hidGlocERUcnlNaXBwT3JzK25iekRLMVY5dWpUUGltZUlTeGsyYk9ETXoKSm8yUTI4YUpGVlJDWFh6TmhaRmFrNFBsaXBQM3BRSm01VEM1Q0RQRzVTQ1NvMTFqeEFYYnlGYzF3SzRIb1lJUwo5QUpEdDRub1lYU1dNUDZubFhiaDZPS2pwcHphNmM1Z0gyM21lZ1ltQ1VLejdDWVljdlRvZjhJN3pkVmxPOFZxCkF0V2hGT3R5cWlpVThncTdBYnBxajE3cG1RYTEyMVpVRGFTdDdVZFgxeTlWbndKR0l1ZU4rRDVBOFdhUmMvZkkKVzl4RDg3V0hJUHMyUG1QK05GdEYxWHovbk45OGppMkw0SGR1dHd2WFFvQ29ESkhLMmVJNWtYODFDVDlJSHFsdApVRUtKR0RNYm04TFhGcFFhTVRuOElaSFpIelk3RG1IOC84aGx6UWlWZkh4QTRKQUJBb0lCQVFDdDJzZUMweUZGCjNpdDdlL1luUnN1NkN4MWs3cCtRQzBZelk4UVR6NE5HYXdzdkVDSWdDcEgvNk9XV1REZXYwYjkvN0VMZkVWRTkKNHJNWGRIMTAvUnRINTM0SEVGVXFzcWN4REE2eXIxVWR2U1lKZ3N4a1RERGNUemRuR2FSYlVNUHg1bWRxck5wZQo2a1k0WkcvVE1QT0J6NmsyLzdhVWNCdGxLUVlJbmZPQ2pNeXpkZXlSSjdJRDFUQXJaUFU3dU9CVFluejhOMWRwCit0cC9iY2o2Z29tU0NaMFBvc3o3NHJUQ3JIOUl2dk1IS0E3ZnVDYnFGeUFEbkJMZEppTWx5bWZhY3VWRCtTeUwKWUVpUFZaMU11OHNlMWwyRjZISHV1eCt2WGtROUdoMzYxMVh2UkR1dncxemlrZTRXQTdocVRQZkZJZy9jRW1GNwpaNEw3MldlOFJyamhBb0lCQVFDc2xpRzNMVGhmWjIzUFozUk1sTll0Y2duSXRDRng3WFliS1l6TnQ1WWhqRFNvCkUwWloxQm5LVDJ2M09YblRpc1RoaTM0SEt5eEJxK3gwSVc3MVhTc1hyWTZLckJsU0FOa2ttYmI3NHZGTnNvYncKamxkcWlHV2dWVEx0UlJuNG1QWUpBNWkrbkUxQ096d1owdHk2TWVDemlSTWpZbGxxZENzTllMSzNEWHkxaWtPKwpnNUsrcXErNzd2NUxTeXVvZ2J4aVIvZ2hVNTV3allHRFJkeTgwY1YvM1hxTngzYUh6bkV1WkV5WXM1dkxpM1dhCnhHajYxM0E4OFlmVXZkYjZ1UTh4ZFVvbHI4NFNjc3VGNHVSNUIvOWxmSUhaMGlkcFV2MFpFNHBmQTR5amhPOSsKMkFRUjROTnkwejJCeFNNZHVJbndxejJueGZmS3JzTVM1YmF5dGFDQkFvSUJBRnp0TWZNVmt4VmJXWGFabmNzRQpwbVI5Q0dzb3VSVXZVWWlxYk9ZQjV6TStpQzNSdTh2UW1wVmxFVUt5M3BrVnpmdzhkc253NGJIb2VMMnl3RlJGCjdjMFRTV1BSTVJTdmhYcEw3WmRJN0lBRzJFd0JJK3NBWnFWN21pdDdvMFJEK1ZoVlJUWFp5cWN0SmZlQ2g5c2sKc1NQVHNhajZLY2RSM1BMSGFMZzJaVENFdmUyMnZJb2g0NTcwMXRoN0VER3A4ZzNmK05wL1lqUDlwOGl4RDlvRwo2QzJ0QWN5WHdtVm9taUhzUGVUT1creVpWc255RHFyVlRZRmdiUnpVQTdseFpPMTR1RjhLMHVwMHZwUU91Uk9JCjFWdFlUWWtENDlJdEp4Ui9tSTNvWmRuc083eTJoZ2krcmVsVkF5TzFQVjlrWUpONFQyM2NUVXErMjE1dXFHb1UKaTBFQ2dnRUFaNGtaYU42RDl4Z0JWRzluNFpsWWM2TDZJNkdNZnVCSi9qbUs4czYwRGlRaVlzSk5iZzVEK281eQo0cmxrVUhmcmJMTldRODZ1bWljZGp2MlBwenJoWXk4SFdFR3VYdmVMVE4yNlhKbmswUXZNei90VWplQ050d1hsCnExbk5Ic29FcjV1c3dvelovR1cweEhrdldiUWFiUnBLbE91bllLbVlPa3BNYkd4MjZDR3VTbGg4YkUzUlp4a1YKRE81bm5vdFdERS9JbDVXbWN6Y3cveU9tTE5CYmZ6M0xDOHNoWEkrSWJxQlZJelo4dkR0SnJqTXVGMjJ2TTZCaQpNRXBDOGQ2Yk1yeCtZdVY0NXJCZlVFNnhhYnBXaVBlTW5yUG9XTk1rYXlyQjFBWTVGTS9uTFYxQjg4ZkFraDRQCnhBNFQ1dnlTSkFOVzFaTjU5K21udFdxQmsreEtBUUtDQVFBTXA4dlYzRUhVTmJRdXBJcTdNUWh6VEdhdUJFZ00KR3p0MnU5cUM0MXYwdS8rMW4wSDlocEx6VGo3VTJ0OHZnSXBoNEtPQSt5VXRUeUt5U3BKbDIxeWthMnp1Wi8rTgpLYUpTcGJmM1BYVUtHNGs3V3dwRGQzaTdoeWw0Z0tMaUNSaWx2Y3cvZmhyQm5UNEhmdGFpN20xbHhMOHBDMk1VCkJmdC9NMkljNU1VVk93Tm5xU012dVhGd3ZHVUpHVWQyWjBBTlZUNmZtT3l3MGp0OFd1R1cxMndwNjV6RUJiYk4KMm1zK0tIU0JyNjJrd1BmNXZaVzZ2K1Rvd0IrdWprM1Y5TzRSZ2MwQXFKYWp3ZkdBS3RKbFZONXMraUQ2RHV5cApLc1RjWklrelpkcERNT2ZnWXh5TVJBWDNUV1lyYnd4V3Q5V1RwSkNkaDZjRWVaRENIMWN1RE9FVAotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQ=="
const PRIVATE_REFRESH_KEY = "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlKS0FJQkFBS0NBZ0VBaUdhTUQ2QS9yaWJWUEtVNzdKUmxuejE3V0tsNmVzVVhpcDFjL3RVQVNuK3FYMlhrCi9XVmZrQWloK2QySldoZEdQQkJGaGdiTXVwTGUzR0wyY2F2bHBaT2wwanMvTW9zNVlVU3BZcVpPVnJydGlQb08KL1NVNWMxSmZPaHpEMGJ3L1Joa0l3OGsvLzJEM015TkZIeGh4OVhhY2Z0aTIzaEFzWk5rUjYwcEtjaWpvcFFFSApHMWIwVGhpWitJSncvdkZLWGxHMjRYSys0K1FkcjlCa253TmVxbkJ0VWZHOXpYaUdRaXN2ZEJLdkNwUi9rdi9KCkhTVmtzejZDNEJPUEtuVFV3OGljR24rVzR2UUptam5scjhoclVVcm16bGM0UW02T2UvR2RGcmJNVkdqcE1rQ0gKK0hNQjFTZDBrMWorYVdsUTExblhidUw2TlE4RjZkY3VVallhbE8zTm41ZG83L3NvWmFNV3dYU2xSYjRDWFJ0QgpTd2ZlOFZrdHRVNDhKUkhSemRORmdKcGU1QllHR2RiTm9LR0VmNW1OWkY3VU81YWZnUENMQWk0bnY5aGNrWVI2CkZnd3laU3E5NnNBS09qT2h6eks3QzVTQ0F5T01zZEg0dUJTaTM0djZtVFNYdDhpWGoyWHNsaVBqU3BPaTBBT28KZXRFbnREY2VPcVl5VzM2c0Q0RjRydWJKU1l1NlVGQU9iWlpZWkppckp0Vnc0Z3diNDFQZzdWbUFUSkh6alI0cgpkMWg1Y3FkNlM5WisyOWtnVExQTEJGMGZSTUw1UzhJSld6bUJuc0g5VTdzNlMxeEVueWVPczAySTlXa2VwOVFkCmZpcUJ1S2p1d3FCblZVVEdvOGZHVEprZW8wSmVKMW4xcTE5dUlYWC82TnVsZTNZcFhHTGxwcFVUbWZNQ0F3RUEKQVFLQ0FnQWdwelAyZGExbytvRG53TUtrc3kzVXZqb3VFbnh3c1lnZU5lZlNWWmw1UERERUg3ZCs5ZXEzcDJsbgpVS0tWLzZaZnNLR0VJVktYZzV0NGRQUjhaK05WRFJDUVVJQ2pqL0xQbDBsWmhXaVJtTFJPcTFZMVFka01BM2NxCmlVSlRqbFl6YU1EUlpmYzlJckJxL0pHS2pTYVMxYTlIS29nMGh2aXB0OUZ6VzFpUkZid0Q3RWdRUW5PLzBtSGgKdlJCaDU4K2UzcjhDSDU4VkhVSUNHY2hNek5pM3dxeFpCcDhpZGl6bDRFSys3YzRib1VzZEhNQy9pbmtkOCtRTwowbi9lY1JPU1B3OG54TTJSVFV5VE1ETU5MdFNLSkgwMmZtaklkb0VEcU9hclRsMkNBNDkzNlR6anZGeCs4N0FUCnFpVlZoTkhHakFwbjBFeUhzUzRBT2ZRcjJDbnFpcG1yNjByK3RhdjNoZWNHV2w1UnNDZktNVkREWDYvT2FPMkkKUU5iTXVxekR2SEJhczJaZm43WGdjMXlhYlBmK3BhQVcrbzZuS20xaW1EMkw0VGI1T21XWWtueW11TnlGUGpqUApvMnpRaHJsdDR1TXlHaGd2MnRiYlBHeEh1emRjMUw3MDRyTnZvd0taZFJrdzBMWmp1b1U0QzZuSUJrdWpwNlh4CmV0OHlpSU1Wa3ZSNTNMZE5sWE5RN2tjMWwrWE1NbXptenQ4RVUvWFdqa1UwdFdjZWJvc0NJSzh2N214NUR0YUgKNU5rM1ppNUJsSHRqNFU1eHZoSi9HQVlON0RSMjBtSmhZcDViSWNwd2RqSDl2UW8xOUQwMkVrc2dHc0NWdWxqegpGZ1ljeEZLL3BGdjQ1ZjdBUDFVZDNKNlRvVG95M01xMzNZZXNFZ1AxQXlFaHhuTWFTUUtDQVFFQTgyMi9RZ3FFCmZNbDFZaFpqS2FMWUM5cG9OK2hnaGoydXhvKytpQWpxNHhmbE5jY3FrSGU2ejNGSjFTd2wwS1lISmRyZ2hQUmgKSGNCUkxZLzUyK2gvbVVMcEIyMlRBUGQyclVzNk5PcE1IaHNEYllLYURueWx2ME54Rll0QWI3QnZ3ZFhFZXpGQgpLaU1EMTAyZjc5V3pmblpNcGRpTDZkR0pGWmpqNFJoSlN6bE5lS3B2RndzYzlKek10OTZzUzBvNGNtVzB5TnZKCk5OS0ZSMmhDNUtaMHhYNkV6RG1WazZFdjRXcytnUTkrTThxTTRYQUJiVXZnVGRPakRsUUxpRWlxU1VvU1pwOHEKT0VINVBDQ0kzNndlWStBbW5kTzhEc3J1UGxIRlBnRUJLRld5aTZuN2lSZDRsalhzRHBMb0QyMnNSRXoxN0RZRgptNHljWHI5RmtTeFUzUUtDQVFFQWozSFZRVFBFVEEzU201WGNyZW5GdUNMOWtTekFJOEVjRFY4ODRnMTU5V1ZjCjJ2bWdWYkhxL2VGU01uZVkyTEhFaTRwKzJQTktxdXlFY0FGcEdpOEdRL3RmdnNWUkVvVWZ3WGt2cG5zeFZhdUkKUXRoUnQzVU9mM0ZQYnVtamNxUEJLVnVsZFgvUWtYYURDMk5DTmRRTWRCdWw5ZlRRcXF6RHRVTlZlK0FyMHR0ZApxcmhXbW5zZlB6TTFxNW9INVZQcEZSSkZaSTBWQWIyN3BHZUVVRFNWM2RQbjV5VFVreDZwV3NlM3dxVGRvWlp2Ck9Cd2VoSTBwbkFicUxKMnpoRGFqTXJxYUI5YUMvZ2Z1UnFMS0tFbCsrMmIzV3o1SDROenkzaWtNbTlVNU1tMmwKbGpBS05zbDNzQ3BNL0pFVlUyWUgybHd1REYrT3dsRExENy9EQ2c2VkR3S0NBUUFKRzQ0UVZueG1mdE1aZkdUeApaZHBYZHpCM0J3YTFmeEZPOUluWVpSMEVxaHcxU3VKWXpXSDc2TzB1UUp5WmxkeW1tZTNVaTZBbWtNOTR1TzVNClFBS21KVTY2ckdyWG1tcWlTVEpBVUpQUWZJcEFTcWFnN0NEM2F2cU1KODJkWUNpT1JBVTU1cm5kYmJuekVFQnYKcExzMmZBNmFGZVFHTjRTOWZoN29pUlFVOEd2cG05YlVNZUkvZEs1a0lyeW5oSHRnTEZYN1BkM2xVQXNVaE40YgoxKzUraFNGSzBzeTUzVW9CVVJYaGxrYk9nVUdNSGJpdjhpck9QcURYSkdYYUQzM3ZpQW53TlB4TFpveUFwMmIzCmwyVDdyNk5DUEczSXorYmlCZ1V2TUxKdVkrWnVPMG5oOHpMYnkrQ3RHdW43eWNxc000VHY5WVY0TUdhWlZPYXgKYThzeEFvSUJBRlhza1JxS2dNWXg3WHRITExaOGR2UlMrV0x4MUhKV1paQlpBU1pEZmsrUmxTcVNKd25PRm41WApieDVONTUrOTlJYkZ3akFBcERSNGt1aG8zK1ZRVDhkL0Z4NDJJZGNmS1NPQ2pSbURaOHp5Z0IvU1pqaW5oTFN3ClVpMlZCRlJTWlExNkdVV2w5M0I2OWdwblBhenl4VGJ6ck5rRStjMlN5WFNWemVueklTMGdQQjVjWjN3SHpuTFUKSVEwV3FpNGpzbFh0Nk9WUFlVcjJ0U1RJNFVnT0I4dWwrSjdMd0E0VWFzdTNJSXNXcUsvM1pjM05nalpTUEo5Ngp3T0ZTNGNxTDAvdzZMTFFQT2M0alFBYk4wcHlKVWVnVUNJMStaQjM5Ry9vWnlyUzU1NVllZWZiWjlmUlZnRDFsClNWSnRNY0lRTnhvRTU0eC8zUXJtekl3MWlRWklMOThDZ2dFQkFKdzFoTGFpNkVnd1pwc3VIV011akIxd1J4dlIKakU0bmxNWk1pVFlua3BRbFhqQ3pWdXQ1NzhsQXpESTFITlE4T2UwMG41T2Yvc0c2THpDOFRIOUJSZFpBWktRUAoycGN2RFpMQ3N4ejhIVVdSSlY3SWpqMlNFN1FQWEZQZ1FYQXcydFpnaDRENit3TE0rWHVSc3hjbEh5TWkva3VECmVPY1pmTlc3dzJnYWhMUUQxZnF2V2QrMUM4K2I5K2ZrMXBMMFo5eG9qZlJpbXRTT20xSWRNakhseHA2QTlGOVoKL1BMc2wyd1preDBCQThZTnhCUGYwcHFqNWEwaFhyblhXVENSZ2lLcExCVTZwdUhoZGFaK05GUUI0bHYvQUcyUAo1c1pJaWkvSG1rc0NTanJoejBjdThra2gra2dleU10SUZVWFgrMkpua0NNVkVLa2x3MXMvdjB3Z3JZdz0KLS0tLS1FTkQgUlNBIFBSSVZBVEUgS0VZLS0tLS0="

func TestVerifyAccessToken(t *testing.T) {
	// Mock dependencies
	mockRedis := new(redis.MockRedis)
	client := new(store.MockMongoClient)

	t.Run("VerifyAccessTokenEmpty", func(t *testing.T) {

		// Initialize the service
		service := user.NewUserService(client, mockRedis)

		// Mock logger
		logger := slog.Default()

		// Mock token
		token := jwt.New(jwt.SigningMethodHS256)
		tokenString, _ := token.SignedString([]byte("secret"))

		// Test the service
		_, err := service.VerifyAccessToken(logger, tokenString)
		assert.Error(t, err)
		assert.Equal(t, "public key not found", err.Error())
	})

	// publicKey, err := base64.StdEncoding.DecodeString(public)
	t.Run("VerifyAccessTokenFail DecodeString err", func(t *testing.T) {
		os.Setenv("PUBLIC_ACCESS_KEY", "!!invalid@@base64^^")
		defer os.Unsetenv("PUBLIC_ACCESS_KEY")

		// Initialize the service
		service := user.NewUserService(client, mockRedis)

		// Mock logger
		logger := slog.Default()

		// Test the service
		_, err := service.VerifyAccessToken(logger, "")
		assert.Error(t, err)
	})

	t.Run("VerifyAccessTokenFail Parse error", func(t *testing.T) {
		os.Setenv("PUBLIC_ACCESS_KEY", "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0")
		defer os.Unsetenv("PUBLIC_ACCESS_KEY")

		// Initialize the service
		service := user.NewUserService(client, mockRedis)

		// Mock logger
		logger := slog.Default()

		// Test the service
		_, err := service.VerifyAccessToken(logger, "")
		assert.Error(t, err)
	})

	// ParseRSAPublicKeyFromPEM error
	t.Run("VerifyAccessTokenFail ParseRSAPublicKeyFromPEM error", func(t *testing.T) {
		invalidPEM := "INVALID_PEM_FORMAT"
		encodedKey := base64.StdEncoding.EncodeToString([]byte(invalidPEM))
		os.Setenv("PUBLIC_ACCESS_KEY", encodedKey)

		// Initialize the service
		service := user.NewUserService(client, mockRedis)

		// Mock logger
		logger := slog.Default()

		// Test the service
		_, err := service.VerifyAccessToken(logger, "")
		assert.Error(t, err)
	})

	t.Run("VerifyAccessTokenFail", func(t *testing.T) {
		os.Setenv("PUBLIC_ACCESS_KEY", "")
		defer os.Unsetenv("PUBLIC_ACCESS_KEY")

		// Initialize the service
		service := user.NewUserService(client, mockRedis)

		// Mock logger
		logger := slog.Default()

		// Mock Redis
		mockRedis.On("Exists", mock.Anything, mock.Anything).Return(int64(1), nil)

		// Test the service
		_, err := service.VerifyAccessToken(logger, "")
		assert.Error(t, err)
	})

	t.Run("ParseWithClaimsError", func(t *testing.T) {
		os.Setenv("PUBLIC_ACCESS_KEY", PUBLIC_ACCESS_KEY)
		defer os.Unsetenv("PUBLIC_ACCESS_KEY")

		// Initialize the service
		service := user.NewUserService(client, mockRedis)

		// Mock logger
		logger := slog.Default()

		// Mock Redis
		mockRedis.On("Exists", mock.Anything, mock.Anything).Return(int64(1), nil)
		invalidToken := "invalid.token.signature"

		// Test the service
		_, err := service.VerifyAccessToken(logger, invalidToken)
		assert.Error(t, err)
	})

	t.Run("VerifyAccessTokenSuccess", func(t *testing.T) {
		os.Setenv("PUBLIC_ACCESS_KEY", PUBLIC_ACCESS_KEY)
		defer os.Unsetenv("PUBLIC_ACCESS_KEY")

		// Initialize the service
		service := user.NewUserService(client, mockRedis)

		// Mock logger
		logger := slog.Default()

		// Mock token

		tokenString, err := mockToken()
		if err != nil {
			t.Fatal("Failed to sign token:", err)
		}
		// Mock Redis
		mockRedis.On("Exists", mock.Anything, mock.Anything).Return(int64(1), nil)

		// Test the service
		result, err := service.VerifyAccessToken(logger, tokenString)
		assert.NoError(t, err)

		// result
		assert.NotNil(t, result)

	})

}

func mockToken() (string, error) {
	privateKey, err := base64.StdEncoding.DecodeString(PRIVATE_ACCESS_KEY)
	if err != nil {
		return "", err
	}

	rsa, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))

	claims := &user.RegisteredClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "1234567890",
			Issuer:    os.Getenv("ISSUER"),
			Audience:  jwt.ClaimStrings{},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
		},
	}

	claims.Email = "test@test.com"
	claims.UserName = "test"
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(rsa)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func mockRefreshToken() (string, error) {
	privateKey, err := base64.StdEncoding.DecodeString(PRIVATE_REFRESH_KEY)
	if err != nil {
		return "", err
	}

	rsa, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))

	claims := &user.RegisteredClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "1234567890",
			Issuer:    os.Getenv("ISSUER"),
			Audience:  jwt.ClaimStrings{},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
		},
	}

	claims.Email = "test@test.com"
	claims.UserName = "test"
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(rsa)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func TestRefreshToken(t *testing.T) {
	// Mock dependencies
	mockRedis := new(redis.MockRedis)
	collection := new(store.MockMongoCollection)
	database := new(store.MockMongoDatabase)
	client := new(store.MockMongoClient)

	t.Run("VerifyRefreshTokenEmpty", func(t *testing.T) {

		// Initialize the service
		service := user.NewUserService(client, mockRedis)

		// Mock logger
		logger := slog.Default()

		// Mock token
		token := jwt.New(jwt.SigningMethodHS256)
		tokenString, _ := token.SignedString([]byte("secret"))

		ctx := context.Background()

		// Test the service
		_, err := service.RefreshToken(ctx, logger, tokenString)
		assert.Error(t, err)
		assert.Equal(t, "public key not found", err.Error())
	})

	t.Run("VerifyRefreshTokenFail DecodeString err", func(t *testing.T) {
		os.Setenv("PUBLIC_REFRESH_KEY", "!!invalid@@base64^^")
		defer os.Unsetenv("PUBLIC_REFRESH_KEY")

		// Initialize the service
		service := user.NewUserService(client, mockRedis)

		// Mock logger
		logger := slog.Default()

		//
		ctx := context.Background()

		// Test the service
		_, err := service.RefreshToken(ctx, logger, "")
		assert.Error(t, err)
	})

	t.Run("VerifyRefreshTokenFail Parse error", func(t *testing.T) {
		os.Setenv("PUBLIC_REFRESH_KEY", "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0")
		defer os.Unsetenv("PUBLIC_REFRESH_KEY")

		// Initialize the service
		service := user.NewUserService(client, mockRedis)

		// Mock logger
		logger := slog.Default()

		// Mock token
		token := jwt.New(jwt.SigningMethodHS256)
		tokenString, _ := token.SignedString([]byte("secret"))

		// mock context
		ctx := context.Background()

		// Test the service
		_, err := service.RefreshToken(ctx, logger, tokenString)
		assert.Error(t, err)
	})

	// ParseRSAPublicKeyFromPEM error
	t.Run("VerifyRefreshTokenFail ParseRSAPublicKeyFromPEM error", func(t *testing.T) {
		invalidPEM := "INVALID_PEM_FORMAT"
		encodedKey := base64.StdEncoding.EncodeToString([]byte(invalidPEM))
		os.Setenv("PUBLIC_REFRESH_KEY", encodedKey)

		// Initialize the service
		service := user.NewUserService(client, mockRedis)

		// Mock logger
		logger := slog.Default()

		ctx := context.Background()

		// Test the service
		_, err := service.RefreshToken(ctx, logger, "")
		assert.Error(t, err)
	})

	t.Run("VerifyRefreshTokenFail", func(t *testing.T) {
		os.Setenv("PUBLIC_REFRESH_KEY", "")
		defer os.Unsetenv("PUBLIC_REFRESH_KEY")

		// Initialize the service
		service := user.NewUserService(client, mockRedis)

		// Mock logger
		logger := slog.Default()

		// Mock token
		token := jwt.New(jwt.SigningMethodHS256)
		tokenString, _ := token.SignedString([]byte("secret"))

		// mock context
		ctx := context.Background()

		// Mock Redis
		mockRedis.On("Exists", mock.Anything, mock.Anything).Return(int64(1), nil)

		// Test the service
		_, err := service.RefreshToken(ctx, logger, tokenString)
		assert.Error(t, err)
	})

	t.Run("VerifyRefreshTokenFail > ParseWithClaims error", func(t *testing.T) {
		os.Setenv("PUBLIC_REFRESH_KEY", PUBLIC_REFRESH_KEY)
		os.Setenv("PRIVATE_ACCESS_KEY", PRIVATE_ACCESS_KEY)
		os.Setenv("PRIVATE_REFRESH_KEY", PRIVATE_REFRESH_KEY)

		defer os.Unsetenv("PUBLIC_REFRESH_KEY")
		defer os.Unsetenv("PRIVATE_ACCESS_KEY")
		defer os.Unsetenv("PRIVATE_REFRESH_KEY")

		// Initialize the service
		service := user.NewUserService(client, mockRedis)

		// Mock logger
		logger := slog.Default()

		// Mock token
		tokenString := "invalid.token.signature"

		// Mock Redis

		mockRedis.On("Exists", mock.Anything, mock.Anything).Return(int64(1), nil)
		mockRedis.On("Del", mock.Anything, mock.Anything).Return(nil)
		mockRedis.On("SetEx", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		// db
		ctx := context.TODO()

		// Test the service
		result, err := service.RefreshToken(ctx, logger, tokenString)
		assert.Error(t, err)

		// result
		assert.Nil(t, result)
	})

	t.Run("VerifyRefreshTokenSuccess", func(t *testing.T) {

		os.Setenv("PUBLIC_REFRESH_KEY", PUBLIC_REFRESH_KEY)
		os.Setenv("PRIVATE_ACCESS_KEY", PRIVATE_ACCESS_KEY)
		os.Setenv("PRIVATE_REFRESH_KEY", PRIVATE_REFRESH_KEY)

		defer os.Unsetenv("PUBLIC_REFRESH_KEY")
		defer os.Unsetenv("PRIVATE_ACCESS_KEY")
		defer os.Unsetenv("PRIVATE_REFRESH_KEY")

		// Initialize the service
		service := user.NewUserService(client, mockRedis)

		// Mock logger
		logger := slog.Default()

		// Mock token
		tokenString, err := mockRefreshToken()
		if err != nil {
			t.Fatal(err)
		}

		// Mock Redis

		mockRedis.On("Exists", mock.Anything, mock.Anything).Return(int64(1), nil)
		mockRedis.On("Del", mock.Anything, mock.Anything).Return(nil)
		mockRedis.On("SetEx", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		// db
		ctx := context.TODO()
		mockResult := mongo.NewSingleResultFromDocument(bson.M{"id": "1234567890", "name": "Test User"}, nil, nil)

		collection.On("FindOne", ctx, bson.M{"id": "1234567890"}, mock.Anything).Return(mockResult)
		client.On("Database", "auth").Return(database)
		database.On("Collection", "users").Return(collection)

		// Test the service
		result, err := service.RefreshToken(ctx, logger, tokenString)
		assert.NoError(t, err)

		// result
		assert.NotNil(t, result)

	})

}

func TestGetUser(t *testing.T) {
	mockClient := new(store.MockMongoClient)
	mockDB := new(store.MockMongoDatabase)
	mockCollection := new(store.MockMongoCollection)
	mockRedis := new(redis.MockRedis)

	mockResult := mongo.NewSingleResultFromDocument(bson.M{"id": "123", "name": "Test User"}, nil, nil)
	ctx := context.TODO()

	mockCollection.On("FindOne", ctx, bson.M{"id": "12345"}, mock.Anything).Return(mockResult)

	// Set up mock database and client
	mockClient.On("Database", "auth").Return(mockDB)
	mockDB.On("Collection", "users").Return(mockCollection)

	// Create the service instance with the mocked client
	service := user.NewUserService(mockClient, mockRedis)

	// Create a mock logger (assuming a simple logger is used)
	mockLogger := slog.Default()

	// Call the method being tested
	u, err := service.GetUser(ctx, mockLogger, "12345")

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert that the returned user ID is correct
	assert.Equal(t, "123", u.ID)

	// Assert that all expected method calls were made
	mockClient.AssertExpectations(t)
	mockDB.AssertExpectations(t)
	mockCollection.AssertExpectations(t)
}
