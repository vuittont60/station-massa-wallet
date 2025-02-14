import { WindowSetSize } from '@wailsjs/runtime/runtime';
import { FiArrowRight } from 'react-icons/fi';

import { PromptRequestData } from '../Sign';
import Intl from '@/i18n/i18n';
import { formatStandard, masToken, maskAddress, Unit } from '@/utils';

export function Transaction(props: PromptRequestData) {
  const {
    WalletAddress,
    RecipientAddress,
    OperationType,
    Amount,
    Fees,
    Expiry,
    Nickname,
  } = props;

  WindowSetSize(460, 500);

  return (
    <div className="flex flex-col items-center gap-4 mas-menu-default">
      <div className="flex w-full items-center justify-between">
        <div className="flex flex-col">
          <div className="flex gap-2">
            <p className="mas-menu-active">
              {Intl.t('password-prompt.sign.from')}
            </p>
            <p className="mas-menu-default">{Nickname}</p>
          </div>
          <p className="mas-caption">{maskAddress(WalletAddress)}</p>
        </div>
        <div className="h-8 w-8 rounded-full flex items-center justify-center bg-neutral">
          <FiArrowRight size={24} className="text-primary" />
        </div>
        <div className="flex flex-col">
          <p className="mas-menu-active">{Intl.t('password-prompt.sign.to')}</p>
          <p className="mas-caption">{maskAddress(RecipientAddress)}</p>
        </div>
      </div>

      <hr className="h-0.25 bg-neutral opacity-40 w-full" />

      <div className="flex justify-between w-full">
        <div className="flex flex-col h-fit">
          <p>{Intl.t('password-prompt.sign.operation-type')}</p>
        </div>
        <div className="flex flex-col items-end h-fit">
          <p>{OperationType}</p>
        </div>
      </div>

      <hr className="h-0.25 bg-neutral opacity-40 w-full" />

      <div className="flex flex-col gap-2 w-full">
        <div className="flex w-full items-center justify-between">
          <p>{Intl.t('password-prompt.sign.sending-amount')}</p>
          <p>
            {formatStandard(Amount, Unit.NanoMAS)} {masToken}
          </p>
        </div>
      </div>

      <div className="flex w-full items-center justify-between">
        <p>{Intl.t('password-prompt.sign.fees')}</p>
        <p>
          {formatStandard(Fees, Unit.NanoMAS)} {masToken}
        </p>
      </div>

      <div className="flex w-full items-center justify-between">
        <p>{Intl.t('password-prompt.sign.expiry')} </p>
        <p>{Expiry}</p>
      </div>

      <hr className="h-0.25 bg-neutral opacity-40 w-full" />
    </div>
  );
}
