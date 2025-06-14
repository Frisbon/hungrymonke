/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * @author ENDERZOMBI102 <enderzombi102.end@gmail.com> 2024
 * @description Quick and dirty `lint-openapi` config to better conform to the Prof's requests and style.
 */
const ibmCloudValidationRules = require( '@ibm-cloud/openapi-ruleset' );
const { allowedKeywords, propertyCasingConvention, operationIdCasingConvention } = require( '@ibm-cloud/openapi-ruleset/src/functions' );
const { operations, schemas, parameters } = require( '@ibm-cloud/openapi-ruleset-utilities/src/collections' );

console.log( 'Loaded config from `.spectral.js`' );

export default {
	extends: ibmCloudValidationRules,
	rules: {
		// REASON: using `examples` instead of `example` as the latter is deprecated
		'ibm-schema-keywords': {
			description: 'Disallows the use of certain keywords',
			message: '{{error}}',
			resolved: true,
			given: schemas,
			severity: 'error',
			then: {
				function: allowedKeywords,
				functionOptions: {
					keywordAllowList: [
						'$ref',
						'additionalProperties',
						'allOf',
						'anyOf',
						'default',
						'description',
						'discriminator',
						'enum',
						'example',
						'exclusiveMaximum',
						'exclusiveMinimum',
						'format',
						'items',
						'maximum',
						'maxItems',
						'maxLength',
						'maxProperties',
						'minimum',
						'minItems',
						'minLength',
						'minProperties',
						'multipleOf',
						'not',
						'oneOf',
						'pattern',
						'patternProperties',
						'properties',
						'readOnly',
						'required',
						'title',
						'type',
						'uniqueItems',
						'unevaluatedProperties',
						'writeOnly',
					]
				}
			}
		},
		// REASON: the operation ids given by the prof follow the camelCase style
		'ibm-operationid-casing-convention': {
			description: 'Operation ids must follow camel case',
			message: '{{error}}',
			resolved: true,
			given: operations,
			severity: 'warn',
			then: {
				function: operationIdCasingConvention,
				functionOptions: {
					type: 'camel'
				}
			}
		},
		// REASON: the prof usually wants camelCase... and i hate snake-case :P
		'ibm-property-casing-convention': {
			description: 'Property names must follow camel case',
			message: '{{error}}',
			resolved: true,
			given: schemas,
			severity: 'warn',
			then: {
				function: propertyCasingConvention,
				functionOptions: {
					type: 'camel'
				}
			}
		},
		// REASON: was enabled by the prof... and it isn't enabled by default
		'ibm-property-consistent-name-and-type': 'warn',
		// REASON: prof elevated these from `warning` to `error`
		'ibm-request-and-response-content': 'error',
		'ibm-avoid-repeating-path-parameters': 'error',
		// REASON: they do not matter to the project's evaluation
		'ibm-operationid-naming-convention': 'off',
		'ibm-pagination-style': 'off',
		'ibm-avoid-inline-schemas': 'off',
		'ibm-requestbody-name': 'off',
		'ibm-error-response-schemas': 'off',
		'ibm-parameter-casing-convention': 'off',
	}
};