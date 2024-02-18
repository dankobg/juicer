import { CastleRights } from './types';

export class CastleRightsHelper {
	constructor(public cr: CastleRights) {}

	toString(): string {
		let fenCastle = '';

		if (this.cr === CastleRights.None) {
			fenCastle += '-';
		}

		if ((this.cr & CastleRights.WhiteKingSide) !== 0) {
			fenCastle += 'K';
		}
		if ((this.cr & CastleRights.WhiteQueenSide) !== 0) {
			fenCastle += 'Q';
		}
		if ((this.cr & CastleRights.BlackKingSide) !== 0) {
			fenCastle += 'k';
		}
		if ((this.cr & CastleRights.BlackQueenSide) !== 0) {
			fenCastle += 'q';
		}

		return fenCastle;
	}
}
