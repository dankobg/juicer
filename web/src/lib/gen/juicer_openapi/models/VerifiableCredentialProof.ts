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

import { mapValues } from '../runtime';
/**
 * 
 * @export
 * @interface VerifiableCredentialProof
 */
export interface VerifiableCredentialProof {
    /**
     * 
     * @type {string}
     * @memberof VerifiableCredentialProof
     */
    jwt?: string;
    /**
     * 
     * @type {string}
     * @memberof VerifiableCredentialProof
     */
    proofType?: string;
}

/**
 * Check if a given object implements the VerifiableCredentialProof interface.
 */
export function instanceOfVerifiableCredentialProof(value: object): value is VerifiableCredentialProof {
    return true;
}

export function VerifiableCredentialProofFromJSON(json: any): VerifiableCredentialProof {
    return VerifiableCredentialProofFromJSONTyped(json, false);
}

export function VerifiableCredentialProofFromJSONTyped(json: any, ignoreDiscriminator: boolean): VerifiableCredentialProof {
    if (json == null) {
        return json;
    }
    return {
        
        'jwt': json['jwt'] == null ? undefined : json['jwt'],
        'proofType': json['proof_type'] == null ? undefined : json['proof_type'],
    };
}

export function VerifiableCredentialProofToJSON(json: any): VerifiableCredentialProof {
    return VerifiableCredentialProofToJSONTyped(json, false);
}

export function VerifiableCredentialProofToJSONTyped(value?: VerifiableCredentialProof | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'jwt': value['jwt'],
        'proof_type': value['proofType'],
    };
}

