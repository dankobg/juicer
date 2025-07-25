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
 * @interface ContinueWithSettingsUiFlow
 */
export interface ContinueWithSettingsUiFlow {
    /**
     * The ID of the settings flow
     * @type {string}
     * @memberof ContinueWithSettingsUiFlow
     */
    id: string;
    /**
     * The URL of the settings flow
     * 
     * If this value is set, redirect the user's browser to this URL. This value is typically unset for native clients / API flows.
     * @type {string}
     * @memberof ContinueWithSettingsUiFlow
     */
    url?: string;
}

/**
 * Check if a given object implements the ContinueWithSettingsUiFlow interface.
 */
export function instanceOfContinueWithSettingsUiFlow(value: object): value is ContinueWithSettingsUiFlow {
    if (!('id' in value) || value['id'] === undefined) return false;
    return true;
}

export function ContinueWithSettingsUiFlowFromJSON(json: any): ContinueWithSettingsUiFlow {
    return ContinueWithSettingsUiFlowFromJSONTyped(json, false);
}

export function ContinueWithSettingsUiFlowFromJSONTyped(json: any, ignoreDiscriminator: boolean): ContinueWithSettingsUiFlow {
    if (json == null) {
        return json;
    }
    return {
        
        'id': json['id'],
        'url': json['url'] == null ? undefined : json['url'],
    };
}

export function ContinueWithSettingsUiFlowToJSON(json: any): ContinueWithSettingsUiFlow {
    return ContinueWithSettingsUiFlowToJSONTyped(json, false);
}

export function ContinueWithSettingsUiFlowToJSONTyped(value?: ContinueWithSettingsUiFlow | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'id': value['id'],
        'url': value['url'],
    };
}

