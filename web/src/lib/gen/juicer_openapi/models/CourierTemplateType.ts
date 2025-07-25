/* tslint:disable */
/* eslint-disable */
/**
 * Juicer schema
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 1.0.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


/**
 * Template type
 * @export
 */
export const CourierTemplateType = {
    RecoveryInvalid: 'recovery_invalid',
    RecoveryValid: 'recovery_valid',
    RecoveryCodeInvalid: 'recovery_code_invalid',
    RecoveryCodeValid: 'recovery_code_valid',
    VerificationInvalid: 'verification_invalid',
    VerificationValid: 'verification_valid',
    VerificationCodeInvalid: 'verification_code_invalid',
    VerificationCodeValid: 'verification_code_valid',
    Stub: 'stub',
    LoginCodeValid: 'login_code_valid',
    RegistrationCodeValid: 'registration_code_valid'
} as const;
export type CourierTemplateType = typeof CourierTemplateType[keyof typeof CourierTemplateType];


export function instanceOfCourierTemplateType(value: any): boolean {
    for (const key in CourierTemplateType) {
        if (Object.prototype.hasOwnProperty.call(CourierTemplateType, key)) {
            if (CourierTemplateType[key as keyof typeof CourierTemplateType] === value) {
                return true;
            }
        }
    }
    return false;
}

export function CourierTemplateTypeFromJSON(json: any): CourierTemplateType {
    return CourierTemplateTypeFromJSONTyped(json, false);
}

export function CourierTemplateTypeFromJSONTyped(json: any, ignoreDiscriminator: boolean): CourierTemplateType {
    return json as CourierTemplateType;
}

export function CourierTemplateTypeToJSON(value?: CourierTemplateType | null): any {
    return value as any;
}

export function CourierTemplateTypeToJSONTyped(value: any, ignoreDiscriminator: boolean): CourierTemplateType {
    return value as CourierTemplateType;
}

